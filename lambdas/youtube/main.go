package main

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
	"time"
	"tony-g/internal/gemini"
	"tony-g/internal/googlesearch"
	"tony-g/internal/googlesheets"
	"tony-g/internal/spotify"
	"tony-g/internal/ssm"
	"tony-g/internal/youtube"

	"github.com/aws/aws-lambda-go/lambda"
)

type Evt struct {
	VideoIds []string `json:"year"`
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
	allVideos := yt.LoadAllPlaylistItems()
	// it's another loop to filter these now, rather than lower down
	// but it's helpful to filter them first
	// to know how many videos were loaded vs how many exist in google sheets
	reviewVideos := youtube.GetReviewVideos(allVideos)

	fmt.Printf("Loaded %d youtube videos\n", len(reviewVideos))
	if len(reviewVideos) == 0 {
		return
	}

	gs := googlesheets.NewClient(googlesheets.Secrets{
		Email:      paramClient.GoogleClientEmail.Value,
		PrivateKey: paramClient.GooglePrivateKey.Value,
	})

	prevVideos := gs.GetYoutubeVideos()
	fmt.Printf("Loaded %d scraped youtube videos from google sheets\n", len(prevVideos))

	prevVideoMap := map[string]bool{}
	for _, v := range prevVideos {
		prevVideoMap[v.Id] = true
	}

	nextVideos := []youtube.PlaylistItem{}
	for _, v := range reviewVideos {
		if prevVideoMap[v.Snippet.ResourceId.VideoId] {
			continue
		}
		if len(evt.VideoIds) > 0 {
			if !slices.Contains(evt.VideoIds, v.Snippet.ResourceId.VideoId) {
				continue
			}
		}

		nextVideos = append(nextVideos, v)
	}
	fmt.Printf("%d Youtube Videos to pull tracks from\n", len(nextVideos))
	if len(nextVideos) == 0 {
		return
	}

	// add oldest videos oldest first
	slices.SortFunc(nextVideos, func(a, z youtube.PlaylistItem) int {
		if a.Snippet.PublishedAt < z.Snippet.PublishedAt {
			return -1
		}
		if a.Snippet.PublishedAt > z.Snippet.PublishedAt {
			return 1
		}
		return 0
	})

	gem := gemini.NewClient(
		paramClient.GeminiApiKey.Value,
	)

	nextTrackRows := []googlesheets.YoutubeTrackRow{}
	nextVideoRows := []googlesheets.YoutubeVideoRow{}
	upper := int(math.Min(float64(len(nextVideos)), 5)) // max 5 videos
	nextVideos = nextVideos[0:upper]
	for i, v := range nextVideos {
		fmt.Printf("Getting tracks from description %d/%d\r", i+1, len(nextVideos))
		nextTracks := gem.ParseYoutubeDescription(v.Snippet.Description)

		nv := googlesheets.YoutubeVideoRow{
			Id:          v.Snippet.ResourceId.VideoId,
			Title:       v.Snippet.Title,
			PublishedAt: v.Snippet.PublishedAt,
			TotalTracks: len(nextTracks),
			FoundTracks: 0,
			AddedAt:     timestamp,
		}
		nextVideoRows = append(nextVideoRows, nv)

		for _, t := range nextTracks {
			r := googlesheets.YoutubeTrackRow{
				Title:            t.Title,
				Artist:           t.Artist,
				Source:           "",
				FoundTrackInfo:   "",
				SpotifyUrl:       "",
				Link:             t.Url,
				VideoId:          v.Snippet.ResourceId.VideoId,
				VideoPublishDate: v.Snippet.PublishedAt,
				AddedAt:          timestamp,
			}

			nextTrackRows = append(nextTrackRows, r)
		}
	}
	fmt.Printf("Gemini found %d tracks in %d video descriptions\n", len(nextTrackRows), len(nextVideos))
	if len(nextTrackRows) == 0 {
		return
	}

	spc := spotify.NewClient(spotify.Secrets{
		ClientId:     paramClient.SpotifyClientId.Value,
		ClientSecret: paramClient.SpotifyClientSecret.Value,
		RefreshToken: paramClient.SpotifyRefreshToken.Value,
	})
	gcs := googlesearch.NewClient(googlesearch.Secrets{
		ApiKey: paramClient.GoogleSearchApiKey.Value,
		Cx:     paramClient.GoogleSearchCx.Value,
	})

	toAddByYear := map[int][]string{}
	foundMap := map[string]int{}
	totalFound := 0
	numSearches := 1
	for i, t := range nextTrackRows {
		fmt.Printf("finding track %d/%d\r", i+1, len(nextTrackRows))

		videoDate, err := time.Parse(time.RFC3339, t.VideoPublishDate)
		year := -1
		if err == nil {
			year = videoDate.Year()
		}

		res := spc.FindTrack(spotify.FindTrackInput{
			Title:  t.Title,
			Artist: t.Artist,
		})
		if len(res) > 0 {
			nextTrackRows[i].Source = "Spotify Search"
			nextTrackRows[i].FoundTrackInfo = spotify.GetTrackInfo(res[0])
			nextTrackRows[i].SpotifyUrl = res[0].ExternalUrls.Spotify
			toAddByYear[year] = append(toAddByYear[year], res[0].Uri)
			foundMap[t.VideoId]++
			totalFound++
			continue
		}

		fmt.Printf("\nGoogle Search %d / 100(quotaLimit): \"%s by %s\"\n", numSearches, t.Title, t.Artist)
		numSearches++

		res2 := gcs.FindSpotifyTrack(googlesearch.FindTrackInput{
			Title:  t.Title,
			Artist: t.Artist,
		})
		if len(res2) == 0 {
			continue
		}

		uri, ok := spotify.LinkToTrackUri(res2[0].Link)
		if !ok {
			continue
		}

		nextTrackRows[i].Source = "Google Search"
		nextTrackRows[i].FoundTrackInfo = res2[0].Title
		nextTrackRows[i].SpotifyUrl = res2[0].Link
		toAddByYear[year] = append(toAddByYear[year], uri)
		foundMap[t.VideoId]++
		totalFound++
	}

	fmt.Printf("Found %d / %d tracks\n", totalFound, len(nextTrackRows))

	myPlaylists := spc.GetMyPlaylists()
	fmt.Printf("Loaded %d playlists\n", len(myPlaylists))
	playlistsByYear := map[int]spotify.SpotifyPlaylist{}
	for _, p := range myPlaylists {
		if strings.HasPrefix(p.Name, spotify.YoutubePlaylistPrefix) {
			yearStr := strings.TrimPrefix(p.Name, spotify.YoutubePlaylistPrefix)
			year, err := strconv.Atoi(yearStr)
			if err == nil {
				playlistsByYear[year] = p
			}
		}
	}
	fmt.Printf("Found %d Melon (Deluxe) playlists\n", len(playlistsByYear))

	yearsSorted := []int{}
	for year := range toAddByYear {
		yearsSorted = append(yearsSorted, year)
	}
	slices.Sort(yearsSorted) // sorts in ascending, which is what I want
	// so that spotify will order playlists from oldest to newest
	for _, year := range yearsSorted {
		uris := toAddByYear[year]

		playlist, ok := playlistsByYear[year]
		fmt.Printf("spotify playlist for %d exists: %t\n", year, ok)

		newTracks := []string{}
		if !ok {
			fmt.Printf("creating spotify playlist: %d\n", year)
			playlistName := spotify.YoutubePlaylistPrefix + strconv.Itoa(year)
			playlist = spc.CreatePlaylist(playlistName)
			newTracks = uris
		} else {
			currentTracks := spc.GetPlaylistItems(playlist.Id)
			fmt.Printf("loaded %d tracks for playlist: %d\n", len(currentTracks), year)

			currentTrackMap := map[string]bool{}
			for _, t := range currentTracks {
				currentTrackMap[t.Track.Uri] = true
			}
			for _, uri := range uris {
				if !currentTrackMap[uri] {
					newTracks = append(newTracks, uri)
				}
			}
		}

		fmt.Printf("adding %d tracks to playlist %s\n", len(newTracks), playlist.Name)
		spc.AddPlaylistItems(playlist.Id, newTracks)
	}

	fmt.Printf("Adding %d track rows to google sheets\n", len(nextTrackRows))
	gs.AddYoutubeTracks(nextTrackRows)
	for i, v := range nextVideoRows {
		nextVideoRows[i].FoundTracks = foundMap[v.Id]
	}

	fmt.Printf("Adding %d video rows to google sheets\n", len(nextVideoRows))
	gs.AddYoutubeVideos(nextVideoRows)
}

func main() {
	lambda.Start(handleLambdaEvent)
}
