package youtube

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/h2non/gock"
	"github.com/joho/godotenv"
)

func TestYoutube(t *testing.T) {
	t.Run("it can create a new youtube client", func(t *testing.T) {
		t.Skip("its messing with the youtube api key")
		testApiKey := "_test_youtubeApiKey"

		yt := NewClient(YtClientConfig{
			ApiKey: testApiKey,
		})

		if yt.apiKey == "" {
			t.Errorf("apiKey not set on Youtube Client")
		}
	})

	t.Run("can load all playlist items", func(t *testing.T) {
		t.Skip("skip test calling YoutubeAPI")

		// Load actual Youtube API Key
		err := godotenv.Load("../../.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		apiKey := os.Getenv("YOUTUBE_API_KEY")
		yt := NewClient(YtClientConfig{
			ApiKey: apiKey,
		})

		items := yt.LoadAllPlaylistItems("PLShqBUh4XbMUKR87gh9ZxUbsYkQoZd7aF")

		if len(items) == 0 {
			t.Errorf("Failed to load playlist items")
		}

		b, err := json.MarshalIndent(items, "", "	")
		if err != nil {
			panic(err)
		}

		err = os.WriteFile("../../data/youtube-videos.json", b, 0666)
		if err != nil {
			panic(err)
		}
	})

	// TODO: mock Youtube HTTP response https://pkg.go.dev/net/http/httptest
	t.Run("makes correctly formatted API call", func(t *testing.T) {
		defer gock.Off()
		// gock.Observe(gock.DumpRequest)

		testApiKey := "_test_youtubeApiKey"
		testPageToken := "_test_pageToken"

		gock.New("https://www.googleapis.com").
			Get("/youtube/v3/playlistItems").
			MatchParam("maxResults", "50").
			MatchParam("playlistId", "PLP4CSgl7K7or84AAhr7zlLNpghEnKWu2c").
			MatchParam("part", "snippet,status").
			MatchParam("key", testApiKey).
			Reply(200).
			JSON(map[string]any{
				"nextPageToken": testPageToken,
				"items": []map[string]any{
					{
						"snippet": map[string]any{
							"resourceId": map[string]any{
								"videoId": "_test_id1",
							},
						},
					},
				},
			})
		// 2nd
		gock.New("https://www.googleapis.com").
			Get("/youtube/v3/playlistItems").
			MatchParam("maxResults", "50").
			MatchParam("playlistId", "PLP4CSgl7K7or84AAhr7zlLNpghEnKWu2c").
			MatchParam("part", "snippet,status").
			MatchParam("key", testApiKey).
			MatchParam("pageToken", testPageToken).
			Reply(200).
			JSON(map[string]any{
				"nextPageToken": "",
				"items": []map[string]any{
					{
						"snippet": map[string]any{
							"resourceId": map[string]any{
								"videoId": "_test_id2",
							},
						},
					},
				},
			})

		yt := NewClient(YtClientConfig{
			ApiKey: testApiKey,
		})

		items := yt.LoadAllPlaylistItems(TonysWeeklyPlaylistId)

		if len(items) != 2 {
			t.Errorf("Expected to load two playlist items received %d", len(items))
		}

		if items[0].Snippet.ResourceId.VideoId != "_test_id1" {
			t.Errorf("Expected test playlist item 1 to have Id _test_id1. Received %s", items[0].Snippet.ResourceId.VideoId)
		}
		if items[1].Snippet.ResourceId.VideoId != "_test_id2" {
			t.Errorf("Expected test playlist item 2 to have Id _test_id2. Received %s", items[1].Snippet.ResourceId.VideoId)
		}
	})

	t.Run("Loads all my playlists", func(t *testing.T) {
		t.Skip("skip test calling YoutubeAPI")

		// Load actual Youtube API Key
		err := godotenv.Load("../../.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		apiKey := os.Getenv("YOUTUBE_API_KEY")
		yt := NewClient(YtClientConfig{
			ApiKey: apiKey,
		})

		items := yt.LoadAllPlaylists()

		if len(items) == 0 {
			t.Errorf("Failed to load playlists")
		}

		b, err := json.MarshalIndent(items, "", "	")
		if err != nil {
			panic(err)
		}

		err = os.WriteFile("../../data/youtube-playlists.json", b, 0666)
		if err != nil {
			panic(err)
		}
	})

	t.Run("set access token works", func(t *testing.T) {
		t.Skip("skip test calling YoutubeAPI")
		err := godotenv.Load("../../.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		yt := NewClient(YtClientConfig{
			ApiKey:       os.Getenv("YOUTUBE_API_KEY"),
			ClientId:     os.Getenv("YOUTUBE_CLIENT_ID"),
			ClientSecret: os.Getenv("YOUTUBE_CLIENT_SECRET"),
			RefreshToken: os.Getenv("YOUTUBE_REFRESH_TOKEN"),
		})

		yt.setAccessToken()

		fmt.Printf("yt access token: %s\n", yt.accessToken)

		if yt.accessToken == "" {
			t.Errorf("Failed to set access token")
		}
	})

	t.Run("CreatePlaylist works", func(t *testing.T) {
		t.Skip("skip test calling YoutubeAPI")
		err := godotenv.Load("../../.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		yt := NewClient(YtClientConfig{
			ApiKey:       os.Getenv("YOUTUBE_API_KEY"),
			ClientId:     os.Getenv("YOUTUBE_CLIENT_ID"),
			ClientSecret: os.Getenv("YOUTUBE_CLIENT_SECRET"),
			RefreshToken: os.Getenv("YOUTUBE_REFRESH_TOKEN"),
		})

		yt.CreatePlaylist("bright", "eyes")
	})

	t.Run("AddPlaylistItems works", func(t *testing.T) {
		t.Skip("skip test calling YoutubeAPI")
		err := godotenv.Load("../../.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		yt := NewClient(YtClientConfig{
			ApiKey:       os.Getenv("YOUTUBE_API_KEY"),
			ClientId:     os.Getenv("YOUTUBE_CLIENT_ID"),
			ClientSecret: os.Getenv("YOUTUBE_CLIENT_SECRET"),
			RefreshToken: os.Getenv("YOUTUBE_REFRESH_TOKEN"),
		})

		yt.AddPlaylistItems("PLShqBUh4XbMUKR87gh9ZxUbsYkQoZd7aF", []string{
			"qIol9hig2G4",
			"qIol9hig2G4",
		})
	})

	t.Run("FindVideo works", func(t *testing.T) {
		t.Skip("skip test calling YoutubeAPI")
		err := godotenv.Load("../../.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		yt := NewClient(YtClientConfig{
			ApiKey:       os.Getenv("YOUTUBE_API_KEY"),
			ClientId:     os.Getenv("YOUTUBE_CLIENT_ID"),
			ClientSecret: os.Getenv("YOUTUBE_CLIENT_SECRET"),
			RefreshToken: os.Getenv("YOUTUBE_REFRESH_TOKEN"),
		})

		results := yt.FindTrack(FindTrackInput{
			Title:  "Grow Wings and Fly",
			Artist: "King Gizzard and the Lizard Wizard",
		})

		if len(results) == 0 {
			log.Fatal("failed to find any search results for Grow Wings and Fly")
		}

		b, err := json.MarshalIndent(results, "", "	")
		if err != nil {
			panic(err)
		}

		err = os.WriteFile("../../data/youtube-search-results.json", b, 0666)
		if err != nil {
			panic(err)
		}
	})

	t.Run("FindVideo works: what about like videos where you just want the audio?", func(t *testing.T) {
		t.Skip("skip test calling YoutubeAPI")
		err := godotenv.Load("../../.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		yt := NewClient(YtClientConfig{
			ApiKey:       os.Getenv("YOUTUBE_API_KEY"),
			ClientId:     os.Getenv("YOUTUBE_CLIENT_ID"),
			ClientSecret: os.Getenv("YOUTUBE_CLIENT_SECRET"),
			RefreshToken: os.Getenv("YOUTUBE_REFRESH_TOKEN"),
		})

		results := yt.FindTrack(FindTrackInput{
			Title:  "Blinding Lights",
			Artist: "The Weeknd",
		})

		if len(results) == 0 {
			log.Fatal("failed to find any search results for Blinding Lights")
		}

		b, err := json.MarshalIndent(results, "", "	")
		if err != nil {
			panic(err)
		}

		err = os.WriteFile("../../data/youtube-search-results.json", b, 0666)
		if err != nil {
			panic(err)
		}
	})

	// findtrack was giving a playlist
	// fixed with type=video query param
	t.Run("FindVideo works: what is this found resource?", func(t *testing.T) {
		t.Skip("skip test calling YoutubeAPI")
		err := godotenv.Load("../../.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		yt := NewClient(YtClientConfig{
			ApiKey:       os.Getenv("YOUTUBE_API_KEY"),
			ClientId:     os.Getenv("YOUTUBE_CLIENT_ID"),
			ClientSecret: os.Getenv("YOUTUBE_CLIENT_SECRET"),
			RefreshToken: os.Getenv("YOUTUBE_REFRESH_TOKEN"),
		})

		results := yt.FindTrack(FindTrackInput{
			Title:  "R.I.P. Meme ft. FrankJavCee",
			Artist: "Hot Dad",
		})

		if len(results) == 0 {
			log.Fatal("failed to find any search results for R.I.P. Meme ft. FrankJavCee")
		}

		b, err := json.MarshalIndent(results, "", "	")
		if err != nil {
			panic(err)
		}

		err = os.WriteFile("../../data/youtube-search-results.json", b, 0666)
		if err != nil {
			panic(err)
		}
	})

	t.Run("FindVideo works: moses!", func(t *testing.T) {
		t.Skip("skip test calling YoutubeAPI")
		err := godotenv.Load("../../.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		yt := NewClient(YtClientConfig{
			ApiKey:       os.Getenv("YOUTUBE_API_KEY"),
			ClientId:     os.Getenv("YOUTUBE_CLIENT_ID"),
			ClientSecret: os.Getenv("YOUTUBE_CLIENT_SECRET"),
			RefreshToken: os.Getenv("YOUTUBE_REFRESH_TOKEN"),
		})

		// moses is an issue
		// can't find even with just single artist
		// I think I should try the videoId from description but I need to be a little smart
		results := yt.FindTrack(FindTrackInput{
			Title:  "O Mistress Mine",
			Artist: "Moses Sumney",
			// Artist: "Moses Sumney, Michael Thurber, Twelfth Night Cast",
		})

		if len(results) == 0 {
			log.Fatal("failed to find any search results for O Mistress Mine")
		}

		b, err := json.MarshalIndent(results, "", "	")
		if err != nil {
			panic(err)
		}

		err = os.WriteFile("../../data/youtube-search-results.json", b, 0666)
		if err != nil {
			panic(err)
		}
	})

	t.Run("Get Videos works", func(t *testing.T) {
		t.Skip("skip test calling YoutubeAPI")
		err := godotenv.Load("../../.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		yt := NewClient(YtClientConfig{
			ApiKey:       os.Getenv("YOUTUBE_API_KEY"),
			ClientId:     os.Getenv("YOUTUBE_CLIENT_ID"),
			ClientSecret: os.Getenv("YOUTUBE_CLIENT_SECRET"),
			RefreshToken: os.Getenv("YOUTUBE_REFRESH_TOKEN"),
		})

		videoIds := []string{
			"lesdjOSFLLY",
			// "upkEBCIZOZA", // private video is not found. That's good.
		}

		results := yt.GetVideosById(videoIds)

		if len(results) == 0 {
			log.Fatalf("Failed to find video lesdjOSFLLY")
		}

		b, err := json.MarshalIndent(results, "", "	")
		if err != nil {
			panic(err)
		}
		err = os.WriteFile("../../data/youtube-getvideo-results.json", b, 0666)
		if err != nil {
			panic(err)
		}
	})
}
