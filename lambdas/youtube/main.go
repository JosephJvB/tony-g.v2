package main

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
	"time"
	"tony-g/internal/gemini"
	"tony-g/internal/googlesheets"
	"tony-g/internal/ssm"
	"tony-g/internal/youtube"

	"github.com/joho/godotenv"
)

type Evt struct {
	VideoIds []string `json:"videoIds"`
}

func handleLambdaEvent(evt Evt) {
	now := time.Now()
	timestamp := now.Format(time.RFC3339)

	paramClient := ssm.NewClient()
	paramClient.LoadParameterValues()

	yt := youtube.NewClient(youtube.YtClientConfig{
		ApiKey:       paramClient.YoutubeApiKey.Value,
		ClientId:     paramClient.YoutubeClientId.Value,
		ClientSecret: paramClient.YoutubeClientSecret.Value,
		RefreshToken: paramClient.YoutubeRefreshToken.Value,
	})
	allTonysVideos := yt.LoadAllPlaylistItems(youtube.TonysWeeklyPlaylistId)
	// remove vid which aren't weekly track reviews
	tonysReviewVideos := youtube.GetReviewVideos(allTonysVideos)

	fmt.Printf("Loaded %d tonys review videos\n", len(tonysReviewVideos))
	if len(tonysReviewVideos) == 0 {
		return
	}

	// sort oldest vids to newest so playlist order is nicely nicely
	slices.SortFunc(tonysReviewVideos, func(a, z youtube.PlaylistItem) int {
		if a.Snippet.PublishedAt < z.Snippet.PublishedAt {
			return -1
		}
		if a.Snippet.PublishedAt > z.Snippet.PublishedAt {
			return 1
		}
		return 0
	})

	gs := googlesheets.NewClient(googlesheets.Secrets{
		Email:      paramClient.GoogleClientEmail.Value,
		PrivateKey: paramClient.GooglePrivateKey.Value,
	})

	prevReviewVideos := gs.GetTonysVideos()
	fmt.Printf("Loaded %d scraped youtube videos from google sheets\n", len(prevReviewVideos))

	prevVideoMap := map[string]bool{}
	for _, v := range prevReviewVideos {
		prevVideoMap[v.Id] = true
	}

	nextVideos := []youtube.PlaylistItem{}
	for _, v := range tonysReviewVideos {
		if prevVideoMap[v.Snippet.ResourceId.VideoId] {
			continue
		}
		// can use event.VideoIds to process a specific set of videos
		if len(evt.VideoIds) > 0 {
			if !slices.Contains(evt.VideoIds, v.Snippet.ResourceId.VideoId) {
				continue
			}
		}

		nextVideos = append(nextVideos, v)
	}
	fmt.Printf("%d Review Videos to pull tracks from\n", len(nextVideos))
	if len(nextVideos) == 0 {
		return
	}

	gem := gemini.NewClient(
		paramClient.GeminiApiKey.Value,
	)

	nextTrackRows := []googlesheets.FoundTrackRow{}
	nextVideoRows := []googlesheets.TonyVideoRow{}
	upper := int(math.Min(float64(len(nextVideos)), 1)) // max 5 videos
	nextVideos = nextVideos[0:upper]
	for i, v := range nextVideos {
		fmt.Printf("Getting tracks from description %d/%d\r", i+1, len(nextVideos))
		nextTracks := gem.ParseYoutubeDescription(v.Snippet.Description)

		nv := googlesheets.TonyVideoRow{
			Id:          v.Snippet.ResourceId.VideoId,
			Title:       v.Snippet.Title,
			PublishedAt: v.Snippet.PublishedAt,
			TotalTracks: len(nextTracks),
			FoundTracks: 0,
			AddedAt:     timestamp,
		}
		nextVideoRows = append(nextVideoRows, nv)

		videoDate, err := time.Parse(time.RFC3339, v.Snippet.PublishedAt)
		year := -1
		if err == nil {
			year = videoDate.Year()
		}

		for _, t := range nextTracks {
			r := googlesheets.FoundTrackRow{
				Title:                  t.Title,
				Artist:                 t.Artist,
				FoundTrackInfo:         "",
				TrackVideoId:           "",
				Link:                   t.Url,
				ReviewVideoId:          v.Snippet.ResourceId.VideoId,
				ReviewVideoPublishDate: v.Snippet.PublishedAt,
				AddedAt:                timestamp,
				Playlist:               strconv.Itoa(year),
			}

			nextTrackRows = append(nextTrackRows, r)
		}
	}
	fmt.Printf("Gemini found %d tracks in %d video descriptions\n", len(nextTrackRows), len(nextVideos))
	if len(nextTrackRows) == 0 {
		return
	}

	// year -> youtube video id
	// used for final youtube.AddPlaylistItems
	toAddByYear := map[string][]string{}
	// just for counting
	foundMap := map[string]int{}
	totalFound := 0
	for i, t := range nextTrackRows {
		fmt.Printf("finding track %d/%d\r", i+1, len(nextTrackRows))

		// have already found an issue where the youtube video from link in description is private
		// so addPlaylistItem fails
		// TODO: verify these are valid ID's before moving on
		// ie: get video by id. Check privacy status
		// https://developers.google.com/youtube/v3/docs/videos/list
		// if idFromLink != "" {
		// 	nextTrackRows[i].TrackVideoId = idFromLink
		// 	nextTrackRows[i].FoundTrackInfo = "id from link"
		// 	toAddByYear[t.Playlist] = append(toAddByYear[t.Playlist], idFromLink)
		// 	foundMap[t.ReviewVideoId]++
		// 	totalFound++
		// 	continue
		// }

		res := yt.FindTrack(youtube.FindTrackInput{
			Artist: t.Artist,
			Title:  t.Title,
		})

		if len(res) > 0 {
			fmt.Println("found track", t.Artist, t.Title, res[0].Id.VideoId)
			nextTrackRows[i].TrackVideoId = res[0].Id.VideoId
			nextTrackRows[i].FoundTrackInfo = res[0].Snippet.Title
			toAddByYear[t.Playlist] = append(toAddByYear[t.Playlist], res[0].Id.VideoId)
			foundMap[t.ReviewVideoId]++
			totalFound++
		}
	}

	fmt.Printf("Found %d / %d tracks\n", totalFound, len(nextTrackRows))

	myPlaylists := yt.LoadAllPlaylists()
	fmt.Printf("Loaded %d playlists\n", len(myPlaylists))
	playlistsByYear := map[string]youtube.Playlist{}
	for _, p := range myPlaylists {
		if strings.HasPrefix(p.Snippet.Title, youtube.PlaylistPrefix) {
			year := strings.TrimPrefix(p.Snippet.Title, youtube.PlaylistPrefix)
			playlistsByYear[year] = p
		}
	}
	fmt.Printf("Found %d Melon playlists\n", len(playlistsByYear))

	for year := range toAddByYear {
		videoIds := toAddByYear[year]

		playlist, ok := playlistsByYear[year]
		fmt.Printf("Youtube playlist for %s exists: %t\n", year, ok)

		newTracks := []string{}
		if !ok {
			fmt.Printf("creating Youtube playlist: %s\n", year)
			playlistName := youtube.PlaylistPrefix + year
			playlist = yt.CreatePlaylist(playlistName, "")
			// addem all
			newTracks = videoIds
		} else {
			currentTracks := yt.LoadAllPlaylistItems(playlist.Id)
			fmt.Printf("loaded %d tracks for playlist: %s\n", len(currentTracks), year)

			// only add not yet added
			currentTrackMap := map[string]bool{}
			for _, t := range currentTracks {
				currentTrackMap[t.Snippet.ResourceId.VideoId] = true
			}
			for _, vid := range videoIds {
				if !currentTrackMap[vid] {
					newTracks = append(newTracks, vid)
				}
			}
		}

		fmt.Printf("adding %d tracks to playlist %s\n", len(newTracks), playlist.Snippet.Title)
		// this method is now adding playlist items 1 at time.
		// google api quota issue probbo
		yt.AddPlaylistItems(playlist.Id, newTracks)
	}

	fmt.Printf("Adding %d track rows to google sheets\n", len(nextTrackRows))
	gs.AddFoundTracks(nextTrackRows)

	for i, v := range nextVideoRows {
		nextVideoRows[i].FoundTracks = foundMap[v.Id]
	}
	fmt.Printf("Adding %d video rows to google sheets\n", len(nextVideoRows))
	gs.AddTonysVideos(nextVideoRows)
}

func main() {
	// lambda.Start(handleLambdaEvent)
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	handleLambdaEvent(Evt{
		VideoIds: []string{},
	})
}
