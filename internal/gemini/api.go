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

type ConfidenceScoreInput struct {
	Index               int    `json:"index"`
	QueryTitle          string `json:"queryTitle"`
	QueryArtist         string `json:"queryArtist"`
	YoutubeVideoTitle   string `json:"youtubeVideoTitle"`
	YoutubeChannelTitle string `json:"youtubeChannelTitle"`
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

// maybe too conservative.
const MAX_CONFIDENCE_INPUTS = 20

func (c *GeminiClient) GenerateConfidenceScores(inputs []ConfidenceScoreInput) []int {
	allScores := []int{}

	totalStart := time.Now()
	for i := 0; i < len(inputs); i += MAX_CONFIDENCE_INPUTS {
		end := min(i+MAX_CONFIDENCE_INPUTS, len(inputs))
		batch := inputs[i:end]

		start := time.Now()
		scores := c.generateConfidenceScores(batch)
		elapsed := time.Since(start)
		fmt.Printf("Generated confidence scores for batch %d-%d in %ssec\n", i, end, elapsed)

		allScores = append(allScores, scores...)
	}
	totalElapsed := time.Since(totalStart)
	fmt.Printf("Generated confidence scores for all %d inputs in %ssec\n", len(inputs), totalElapsed)

	return allScores
}

// this prompt is a lot slower than the other one
// 5-6sec vs 10-15sec
// obv it's lots bigger
// but defs pref accuracy over speed
// I wanna improve the other one too
func (c *GeminiClient) generateConfidenceScores(inputs []ConfidenceScoreInput) []int {
	input := `The following list is the result of multiple Youtube Search API calls to find songs in Youtube.
Your task is to assign a confidence score from 0 to 100 for each item in the list.
The confidence score should indicate how well the Youtube search API result matches the query.
Each item in the list has the following properties:
- index: Each item has an index field. Your output array must have the same length as the input, where output[i] is the score for input[i].
- queryTitle: song title used in YoutubeSearch query
- queryArtist: artist name used in YoutubeSearch query
- youtubeVideoTitle: title of the YouTube video
- youtubeChannelTitle: name of the YouTube channel that uploaded the video

### Scoring guidelines:

- 90-100: Strong confidence. 
	- "youtubeVideoTitle" should include "queryTitle".
	- "youtubeVideoTitle" might also include "queryArtist".
	- "youtubeVideoTitle" might include "official audio", "official video". This should not lower score.
	- "youtubeChannelTitle" should match "queryArtist".
	- "youtubeChannelTitle" might include "- Topic" or "VEVO". This should not lower score.
	- Examples:
		- {"queryTitle":"Goose Snow Cone","queryArtist":"Aimee Mann","youtubeVideoTitle":"Aimee Mann - Goose Snow Cone (Official Audio)","youtubeChannelTitle":"Aimee Mann"} -> 100
		- {"queryTitle":"Real Death","queryArtist":"Mount Eerie","youtubeVideoTitle":"Real Death","youtubeChannelTitle":"Mount Eerie"} -> 100

- 60-89: Good confidence. Accept cosmetic inconsistencies such as:
	- "youtubeChannelTitle" might not match at all. It may be the offical Music Label releasing the song instead of the artist.
	- "youtubeVideoTitle" might have inconsistent guest features. eg: Feat. ft.
	- Examples:
		- {"queryTitle":"Help ft. Wiki & Edan","queryArtist":"Your Old Droog","youtubeVideoTitle":"Your Old Droog - "Help" feat. Wiki and Edan","youtubeChannelTitle":"Your Old Droog"} -> 89 (imperfect guest features)
		- {"queryTitle":"Memories Are Now","queryArtist":"Jesca Hoop","youtubeVideoTitle":"Jesca Hoop - Memories Are Now [OFFICIAL VIDEO]","youtubeChannelTitle":"Sub Pop"} -> 89 (channel is Music Label)

- 20-59: Moderate confidence. Accept more significant inconsistencies such as:
	- spelling or wording matches are not exact but the sentiment is correct.
	- "youtubeChannelTitle" might be a fan YouTube channel instead of the "queryArtist" or official music label.
	- Examples:
		- {"queryTitle":"Vapid Feels Are Vapid","queryArtist":"Clarence Clarity","youtubeVideoTitle":"Vapid Feels Ain't Vapid","youtubeChannelTitle":"Clarence Clarity - Topic"} -> 59 (incorrect title, correct sentiment)
		- {"queryTitle":"It Takes Two","queryArtist":"Mike Will Made-It, Lil Yachty, Carly Rae Jepsen","youtubeVideoTitle":"Mike WiLL Made It, Lil Yachty, Carly Rae Jepsen   It Takes Two","youtubeChannelTitle":"Андрей Гантимуровl"} -> 30 (fan channel)

- 0-19: No confidence.
	- It's clear the YouTube search result is not the song from the Query.
	- Examples:
		- {"queryTitle":"K33p Ur Dr34ms (Suicide Remix)","queryArtist":"DJ Windows 98 (Win Butler)","youtubeVideoTitle":"Лучшие и худшие треки недели. 24 июня Desiigner, Skrillex, Rick Ross(theneddledrop на русском)","youtubeChannelTitle":"he-he he-he-he"} -> 0 (completely wrong)
		- {"queryTitle":"Dress (PJ Harvey Cover)","queryArtist":"Buke & Gase","youtubeVideoTitle":"Buke & Gase - Hiccup","youtubeChannelTitle":"shepritzl"} -> 19 (right artist wrong song)

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

	scores := []int{}
	err = json.Unmarshal([]byte(result.Text()), &scores)
	if err != nil {
		log.Fatalf("GenerateConfidenceScores: Failed to parse response JSON")
	}

	if len(scores) != len(inputs) {
		log.Fatalf("GenerateConfidenceScores: expected %d scores, got %d", len(inputs), len(scores))
	}

	// d, err := json.MarshalIndent(scores, "", "	")
	// if err != nil {
	// 	panic(err)
	// }

	// err = os.WriteFile("../../data/scorrreee.json", d, 0666)
	// if err != nil {
	// 	panic(err)
	// }

	return scores
}

// TODO: improve prompt
// - give examples of double songs
func (c *GeminiClient) ParseYoutubeDescription(description string) []ParsedTrack {
	input := "Return the Best Tracks mentioned in the following text snippet"
	// it was giving me ...meh... tracks
	input += "\nIgnore the tracks in the \"meh\" and \"worst\" sections"
	input += "\nformat \"{artist} - {title}\n{url}\""
	// handle multi track for one artist case
	input += "\nIf title has one or more slash character and there is more than one url, return multiple tracks and split the titles by slash character"
	input += "\n"
	input += description

	start := time.Now()
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
	elapsed := time.Since(start)
	fmt.Printf("Gemini parsed Youtube description in %ssec\n", elapsed)

	parsedTracks := []ParsedTrack{}
	err = json.Unmarshal([]byte(result.Text()), &parsedTracks)
	if err != nil {
		log.Fatalf("ParseYoutubeDescription: Failed to parse response JSON")
	}

	return parsedTracks
}
