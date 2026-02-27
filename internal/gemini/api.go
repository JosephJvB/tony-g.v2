package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"google.golang.org/genai"
)

type ParsedTrack struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Url    string `json:"url"`
}

type ConfidenceScore struct {
	Index               int    `json:"index"`
	Query               string `json:"query"`
	YoutubeVideoTitle   string `json:"youtubeVideoTitle"`
	YoutubeChannelTitle string `json:"youtubeChannelTitle"`
	Score               int    `json:"score"`
}
type GeminiClient struct {
	client genai.Client
	ctx    context.Context
}

const GEMINI_MODEL = "gemini-2.5-flash"

func NewClient(apiKey string) GeminiClient {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatal(err)
	}

	return GeminiClient{
		client: *client,
		ctx:    ctx,
	}
}

func (c *GeminiClient) GenerateConfidenceScores(inputs []ConfidenceScore) []ConfidenceScore {
	input := `Given a list of items where each has a "query" (artist + song title) and a
"youtubeSearchResult" (the title of the top YouTube search result), score
how confident you are (0-100) that the YouTube result is the correct
official music track.

Scoring guidelines:
- 90-100: Result clearly matches the artist and song title (official audio/video)
- 50-89: Result likely matches but has extra info (feat. artists, remix labels, etc.)
- 20-49: Result is ambiguous — could be a cover, live version, or compilation
- 0-19: Result is clearly wrong (different song, reaction video, unrelated content)

Return a JSON array with the same length as the input
where each item is your confidence score relating to the corresponding input item.`

	jsonList, err := json.Marshal(inputs)
	if err != nil {
		panic(err)
	}

	input += "\n\nInput:\n"
	input += string(jsonList)

	// err = os.WriteFile("../../data/confidence-input.txt", []byte(input), 0666)
	// if err != nil {
	// 	panic(err)
	// }

	result, err := c.client.Models.GenerateContent(
		c.ctx,
		GEMINI_MODEL,
		genai.Text(input),
		&genai.GenerateContentConfig{
			ResponseMIMEType: "application/json",
			ResponseSchema: &genai.Schema{
				Type: genai.TypeArray,
				Items: &genai.Schema{
					Type: genai.TypeInteger,
				},
			},
		},
	)

	if err != nil {
		errStr := err.Error()
		// I think this was wifi cutting out. Rm for now
		// if strings.HasPrefix(errStr, "Error 503, Message: The service is currently unavailable") {
		// 	fmt.Println("Gemini 503 error, 10 sec timeout")
		// 	time.Sleep(time.Second * 10)
		// 	return c.GenerateConfidenceScores(inputs)
		// }
		if strings.HasPrefix(errStr, "Error 429, Message: You exceeded your current quota") {
			fmt.Println("Gemini Quota Exceeded, 10 sec timeout")
			time.Sleep(time.Second * 10)
			return c.GenerateConfidenceScores(inputs)
		}

		log.Fatal(err)
	}

	outputs := []int{}
	err = json.Unmarshal([]byte(result.Text()), &outputs)
	if err != nil {
		log.Fatalf("GenerateConfidenceScores: Failed to parse response JSON")
	}

	if len(outputs) != len(inputs) {
		log.Fatalf("GenerateConfidenceScores: expected %d scores, got %d", len(inputs), len(outputs))
	}

	for i, score := range outputs {
		inputs[i].Score = score
	}

	// d, err := json.MarshalIndent(scores, "", "	")
	// if err != nil {
	// 	panic(err)
	// }

	// err = os.WriteFile("../../data/scorrreee.json", d, 0666)
	// if err != nil {
	// 	panic(err)
	// }

	return inputs
}

func (c *GeminiClient) ParseYoutubeDescription(description string) []ParsedTrack {
	input := "Return the Best Tracks mentioned in the following text snippet"
	// it was giving me ...meh... tracks
	input += "\nIgnore the tracks in the \"meh\" and \"worst\" sections"
	input += "\nformat \"{artist} - {title}\n{url}\""
	// handle multi track for one artist case
	input += "\nIf title has one or more slash character and there is more than one url, return multiple tracks and split the titles by slash character"
	input += "\n"
	input += description

	result, err := c.client.Models.GenerateContent(
		c.ctx,
		GEMINI_MODEL,
		genai.Text(input),
		&genai.GenerateContentConfig{
			// Tools: []*genai.Tool{
			// 	{GoogleSearch: &genai.GoogleSearch{}},
			// },
			// can return JSON but not with a google search!
			ResponseMIMEType: "application/json",
			ResponseSchema: &genai.Schema{
				Type: genai.TypeArray,
				Items: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"title":  {Type: genai.TypeString},
						"artist": {Type: genai.TypeString},
						"url":    {Type: genai.TypeString},
					},
					// haven't tried this
					// idk if I need it
					// But Conall is using this property: https://github.com/schwart/Pull-Request-Writer/blob/master/gemini.go
					// Required: []string{"title", "artist"},
				},
			},
		},
	)
	// try avoid rate limit
	// https://ai.google.dev/gemini-api/docs/rate-limits
	// I'm currently on gemini-2.0-flash
	// 15 requests per minute
	// 200 requests per day
	// https://console.cloud.google.com/iam-admin/quotas?authuser=1&inv=1&invt=AbxR6Q&project=tnd-best-tracks&pageState=(%22allQuotasTable%22:(%22f%22:%22%255B%255D%22))
	// even tho the docs say I should be allowed 1500 per day, my cloud console quota is 1000
	if err != nil {
		errStr := err.Error()
		// I think this was wifi cutting out. Rm for now
		// if strings.HasPrefix(errStr, "Error 503, Message: The service is currently unavailable") {
		// 	fmt.Println("Gemini 503 error, 10 sec timeout")
		// 	time.Sleep(time.Second * 10)
		// 	return c.ParseYoutubeDescription(description)
		// }
		if strings.HasPrefix(errStr, "Error 429, Message: You exceeded your current quota") {
			fmt.Println("Gemini Quota Exceeded, 10 sec timeout")
			time.Sleep(time.Second * 10)
			return c.ParseYoutubeDescription(description)
		}

		log.Fatal(err)
	}

	parsedTracks := []ParsedTrack{}
	err = json.Unmarshal([]byte(result.Text()), &parsedTracks)
	if err != nil {
		log.Fatalf("ParseYoutubeDescription: Failed to parse response JSON")
	}

	return parsedTracks
}
