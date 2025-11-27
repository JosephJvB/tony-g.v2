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

	t.Run("can load all youtube items", func(t *testing.T) {
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

		items := yt.LoadAllPlaylistItems()

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

		items := yt.LoadAllPlaylistItems()

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
}
