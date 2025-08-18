package youtube

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const BaseUrl = "https://www.googleapis.com/youtube/v3"
const PlaylistId = "PLP4CSgl7K7or84AAhr7zlLNpghEnKWu2c"
const MyChannelId = "UC7bbCeEhOfxos9EmsvaxNGQ"

type YoutubeClient struct {
	apiKey string
}

type PlaylistItem struct {
	Snippet struct {
		Title               string `json:"title"`
		Description         string `json:"description"`
		PublishedAt         string `json:"publishedAt"`
		VideoOwnerChannelId string `json:"videoOwnerChannelId"`
		ChannelId           string `json:"channelId"`
		ResourceId          struct {
			Kind    string `json:"kind"`
			VideoId string `json:"videoId"`
		} `json:"resourceId"`
	} `json:"snippet"`
	Status struct {
		PrivacyStatus string `json:"privacyStatus"`
	} `json:"status"`
}
type Playlist struct {
	Id      string `json:"id"`
	Snippet struct {
		PublishedAt string `json:"publishedAt"`
		ChannelId   string `json:"channelId"`
		Title       string `json:"title"`
		Description string `json:"description"`
	} `json:"snippet"`
	Status struct {
		PrivacyStatus string `json:"privacyStatus"`
	} `json:"status"`
	ContentDetails struct {
		ItemCount int `json:"itemCount"`
	} `json:"contentDetails"`
}

type ApiResponse[T Playlist | PlaylistItem] struct {
	NextPageToken string `json:"nextPageToken"`
	Items         []T    `json:"items"`
}

func (yt *YoutubeClient) LoadAllPlaylistItems() []PlaylistItem {
	resp := getPlaylistItems(
		yt.apiKey,
		PlaylistId,
		"",
	)

	items := resp.Items
	pageToken := resp.NextPageToken

	for pageToken != "" {
		resp := getPlaylistItems(
			yt.apiKey,
			PlaylistId,
			pageToken,
		)

		items = append(items, resp.Items...)
		pageToken = resp.NextPageToken
	}

	return items
}

func getPlaylistItems(key string, playlistId string, pageToken string) ApiResponse[PlaylistItem] {
	apiUrl := BaseUrl + "/playlistItems"

	queryPart := url.Values{}
	queryPart.Set("maxResults", "50")
	queryPart.Set("playlistId", playlistId)
	queryPart.Set("part", "snippet,status")
	queryPart.Set("key", key)
	if pageToken != "" {
		queryPart.Set("pageToken", pageToken)
	}

	apiUrl += "?" + queryPart.Encode()

	resp, err := http.Get(apiUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		b := new(strings.Builder)
		io.Copy(b, resp.Body)
		log.Print(b.String())
		log.Fatalf("\ngetPlaylistItems failed: \"%s\"", resp.Status)
	}

	responseBody := ApiResponse[PlaylistItem]{}
	json.NewDecoder(resp.Body).Decode(&responseBody)

	return responseBody
}

func (yt *YoutubeClient) LoadAllPlaylists() []Playlist {
	resp := getPlaylists(
		yt.apiKey,
		MyChannelId,
		"",
	)

	items := resp.Items
	pageToken := resp.NextPageToken

	for pageToken != "" {
		resp := getPlaylists(
			yt.apiKey,
			MyChannelId,
			pageToken,
		)

		items = append(items, resp.Items...)
		pageToken = resp.NextPageToken
	}

	return items
}

func getPlaylists(key string, channelId string, pageToken string) ApiResponse[Playlist] {
	apiUrl := BaseUrl + "/playlists"

	queryPart := url.Values{}
	queryPart.Set("maxResults", "50")
	queryPart.Set("channelId", channelId)
	queryPart.Set("part", "snippet,status,contentDetails")
	queryPart.Set("key", key)
	if pageToken != "" {
		queryPart.Set("pageToken", pageToken)
	}

	apiUrl += "?" + queryPart.Encode()

	resp, err := http.Get(apiUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		b := new(strings.Builder)
		io.Copy(b, resp.Body)
		log.Print(b.String())
		log.Fatalf("\ngetPlaylists failed: \"%s\"", resp.Status)
	}

	responseBody := ApiResponse[Playlist]{}
	json.NewDecoder(resp.Body).Decode(&responseBody)

	return responseBody
}

func NewClient(apiKey string) YoutubeClient {
	return YoutubeClient{
		apiKey: apiKey,
	}
}
