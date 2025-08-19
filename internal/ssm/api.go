package ssm

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type Parameter struct {
	Name  string
	Value string
}

type SsmClient struct {
	client              *ssm.Client
	GoogleClientEmail   Parameter
	GooglePrivateKey    Parameter
	SpotifyClientId     Parameter
	SpotifyClientSecret Parameter
	SpotifyRefreshToken Parameter
	GeminiApiKey        Parameter
	GoogleSearchApiKey  Parameter
	GoogleSearchCx      Parameter
	YoutubeApiKey       Parameter
	YoutubeClientId     Parameter
	YoutubeClientSecret Parameter
	YoutubeRefreshToken Parameter
}

func NewClient() SsmClient {
	awsRegion := os.Getenv("AWS_REGION")
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if err != nil {
		panic(err)
	}

	ssm := ssm.NewFromConfig(cfg)

	return SsmClient{
		client: ssm,
		// Google
		GoogleClientEmail: Parameter{
			Name:  "/Tony2Stack/google/client-email",
			Value: "",
		},
		GooglePrivateKey: Parameter{
			Name:  "/Tony2Stack/google/private-key",
			Value: "",
		},
		// Spotify,
		SpotifyClientId: Parameter{
			Name:  "/Tony2Stack/spotify/client-id",
			Value: "",
		},
		SpotifyClientSecret: Parameter{
			Name:  "/Tony2Stack/spotify/secret",
			Value: "",
		},
		SpotifyRefreshToken: Parameter{
			Name:  "/Tony2Stack/spotify/refresh-token",
			Value: "",
		},
		// Gemini
		GeminiApiKey: Parameter{
			Name:  "/Tony2Stack/google/gemini/api-key",
			Value: "",
		},
		// Google Search
		GoogleSearchApiKey: Parameter{
			Name:  "/Tony2Stack/google/search/api-key",
			Value: "",
		},
		GoogleSearchCx: Parameter{
			Name:  "/Tony2Stack/google/search/cx",
			Value: "",
		},
		YoutubeApiKey: Parameter{
			Name:  "/Tony2Stack/youtube/api-key",
			Value: "",
		},
		YoutubeClientId: Parameter{
			Name:  "/Tony2Stack/youtube/client-id",
			Value: "",
		},
		YoutubeClientSecret: Parameter{
			Name:  "/Tony2Stack/youtube/client-secret",
			Value: "",
		},
		YoutubeRefreshToken: Parameter{
			Name:  "/Tony2Stack/youtube/refresh-token",
			Value: "",
		},
	}
}

func (sc *SsmClient) LoadParameterValues() {
	toLoad := []Parameter{
		sc.GoogleClientEmail,
		sc.GooglePrivateKey,
		sc.SpotifyClientId,
		sc.SpotifyClientSecret,
		sc.SpotifyRefreshToken,
		sc.GeminiApiKey,
		sc.GoogleSearchApiKey,
		sc.GoogleSearchCx,
		sc.YoutubeApiKey,
		sc.YoutubeClientId,
		sc.YoutubeClientSecret,
		sc.YoutubeRefreshToken,
	}

	names := []string{}
	for _, p := range toLoad {
		names = append(names, p.Name)
	}

	params := ssm.GetParametersInput{
		Names: names,
	}

	r, err := sc.client.GetParameters(context.TODO(), &params)
	if err != nil {
		panic(err)
	}

	if len(r.InvalidParameters) != 0 {
		log.Fatalf("Invalid Parameters:%v\n", r.InvalidParameters)
	}

	m := map[string]string{}
	for _, p := range r.Parameters {
		m[*p.Name] = *p.Value
	}

	sc.GoogleClientEmail.Value = m[sc.GoogleClientEmail.Name]
	sc.GooglePrivateKey.Value = m[sc.GooglePrivateKey.Name]
	sc.SpotifyClientId.Value = m[sc.SpotifyClientId.Name]
	sc.SpotifyClientSecret.Value = m[sc.SpotifyClientSecret.Name]
	sc.SpotifyRefreshToken.Value = m[sc.SpotifyRefreshToken.Name]
	sc.GeminiApiKey.Value = m[sc.GeminiApiKey.Name]
	sc.GoogleSearchApiKey.Value = m[sc.GoogleSearchApiKey.Name]
	sc.GoogleSearchCx.Value = m[sc.GoogleSearchCx.Name]
	sc.YoutubeApiKey.Value = m[sc.YoutubeApiKey.Name]
	sc.YoutubeClientId.Value = m[sc.YoutubeClientId.Name]
	sc.YoutubeClientSecret.Value = m[sc.YoutubeClientSecret.Name]
	sc.YoutubeRefreshToken.Value = m[sc.YoutubeRefreshToken.Name]
}
