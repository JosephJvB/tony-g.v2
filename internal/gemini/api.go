package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/genai"
)

type ParsedTrack struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Url    string `json:"url"`
}

type ConfidenceScoresInput struct {
	Index               int    `json:"index"`
	Query               string `json:"query"`
	YoutubeSearchResult string `json:"YoutubeSearchResult"`
}
type ConfidenceScoresOutput struct {
	Index int `json:"index"`
	Score int `json:"score"`
}
type GeminiClient struct {
	client genai.Client
	ctx    context.Context
}

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

func (c *GeminiClient) GenerateConfidenceScores(inputs []ConfidenceScoresInput) []ConfidenceScoresOutput {
	input := "the following list is songs that I have searched for via Youtube Search API"
	input += "\nfor each item, generate a confidence score that the query has returned the correct youtube video"
	input += "\nscore should be between 0 and 100"
	input += "\nthe list of scores must be in the same order as the inputs so I can correctly assign scores"
	input += "\n"
	input += "\nlist:\n"

	jsonList, err := json.Marshal(inputs)
	if err != nil {
		panic(err)
	}

	input += string(jsonList)

	err = os.WriteFile("../../data/confidence-input.txt", []byte(input), 0666)
	if err != nil {
		panic(err)
	}

	result, err := c.client.Models.GenerateContent(
		c.ctx,
		"gemini-2.0-flash",
		genai.Text(input),
		&genai.GenerateContentConfig{
			ResponseMIMEType: "application/json",
			ResponseSchema: &genai.Schema{
				Type: genai.TypeArray,
				Items: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"index": {Type: genai.TypeNumber},
						"score": {Type: genai.TypeNumber},
					},
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

	outputs := []ConfidenceScoresOutput{}
	err = json.Unmarshal([]byte(result.Text()), &outputs)
	if err != nil {
		log.Fatalf("GenerateConfidenceScores: Failed to parse response JSON")
	}

	// d, err := json.MarshalIndent(scores, "", "	")
	// if err != nil {
	// 	panic(err)
	// }

	// err = os.WriteFile("../../data/scorrreee.json", d, 0666)
	// if err != nil {
	// 	panic(err)
	// }

	return outputs
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
		"gemini-2.0-flash",
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
