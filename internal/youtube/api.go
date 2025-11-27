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
	apiKey       string
	clientId     string
	clientSecret string
	refreshToken string
	accessToken  string
}

type PlaylistSnippet struct {
	Title               string `json:"title"`
	Description         string `json:"description"`
	PublishedAt         string `json:"publishedAt"`
	VideoOwnerChannelId string `json:"videoOwnerChannelId"`
	ChannelId           string `json:"channelId"`
	ResourceId          struct {
		Kind    string `json:"kind"`
		VideoId string `json:"videoId"`
	} `json:"resourceId"`
}

type PlaylistItem struct {
	Snippet PlaylistSnippet `json:"snippet"`
	Status  struct {
		PrivacyStatus string `json:"privacyStatus"`
	} `json:"status"`
}
type Playlist struct {
	Id      string          `json:"id"`
	Snippet PlaylistSnippet `json:"snippet"`
	Status  struct {
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

func (yt *YoutubeClient) setAccessToken() {
	if yt.accessToken != "" {
		return
	}

	queryPart := url.Values{}
	queryPart.Set("client_id", yt.clientId)
	queryPart.Set("client_secret", yt.clientSecret)
	queryPart.Set("refresh_token", yt.refreshToken)
	queryPart.Set("grant_type", "refresh_token")

	resp, err := http.Post("https://oauth2.googleapis.com/token", "application/x-www-form-urlencoded", strings.NewReader(queryPart.Encode()))
	if err != nil {
		log.Fatalf("Failed to get access token: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		b := new(strings.Builder)
		io.Copy(b, resp.Body)
		log.Print(b.String())
		log.Fatalf("Failed to get access token: %s", resp.Status)
	}

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		log.Fatalf("Failed to decode access token response: %v", err)
	}

	yt.accessToken = tokenResponse.AccessToken
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

func (yt *YoutubeClient) CreatePlaylist(title string, description string) {
	if yt.accessToken == "" {
		yt.setAccessToken()
	}

	apiUrl := BaseUrl + "/playlists"

	queryPart := url.Values{}
	queryPart.Set("part", "snippet,status")

	apiUrl += "?" + queryPart.Encode()

	postData := map[string]any{
		"snippet": map[string]any{
			"title":       title,
			"description": description,
		},
		"status": map[string]any{
			"privacyStatus": "public",
		},
	}
	postBuffer, _ := json.Marshal(postData)
	postString := strings.NewReader(string(postBuffer))
	log.Println(postString)

	req, _ := http.NewRequest("POST", apiUrl, postString)

	authHeaderValue := "Bearer " + yt.accessToken
	req.Header.Set("Authorization", authHeaderValue)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode > 299 {
		b := new(strings.Builder)
		io.Copy(b, resp.Body)
		log.Print(b.String())
		log.Fatalf("\nCreatePlaylist failed: \"%s\"", resp.Status)
	}

	b := new(strings.Builder)
	io.Copy(b, resp.Body)
	log.Print(b.String())
	log.Println("Create playlist sucess")
}

type YtClientConfig struct {
	ApiKey       string
	ClientId     string
	ClientSecret string
	RefreshToken string
}

func NewClient(config YtClientConfig) YoutubeClient {
	return YoutubeClient{
		apiKey:       config.ApiKey,
		clientId:     config.ClientId,
		clientSecret: config.ClientSecret,
		refreshToken: config.RefreshToken,
	}
}
