package googlesheets

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/joho/godotenv"
)

func TestGoogleSheets(t *testing.T) {

	t.Run("append to undefined map key", func(t *testing.T) {
		m := map[int][]string{}

		// m[20] is not set - will append throw?
		m[20] = append(m[20], "123")

		// t.Logf("%v", m)

		if len(m[20]) != 1 {
			t.Error("something went wrong")
		}
	})
	t.Run("can load videos from google sheets", func(t *testing.T) {
		t.Skip("skip test calling real google sheets api")

		err := godotenv.Load("../../.env")
		if err != nil {
			panic(err)
		}

		// .env file does not handle private keys gracefully
		// probably would be better saved to a file than in .env. Oh well.
		invalidKey := os.Getenv("GOOGLE_SHEETS_PRIVATE_KEY")
		fixedKey := strings.ReplaceAll(invalidKey, "__n__", "\n")

		gs := NewClient(Secrets{
			Email:      os.Getenv("GOOGLE_SHEETS_EMAIL"),
			PrivateKey: fixedKey,
		})

		tracks := gs.GetFoundTracks()

		if len(tracks) == 0 {
			t.Errorf("Expected parsed videos to be loaded")
		}

		b, err := json.MarshalIndent(tracks, "", "	")
		if err != nil {
			panic(err)
		}

		err = os.WriteFile("../../data/scraped-tracks.json", b, 0666)
		if err != nil {
			panic(err)
		}
	})

	t.Run("can append tracks to google sheets", func(t *testing.T) {
		t.Skip("skip test calling real google sheets api")

		err := godotenv.Load("../../.env")
		if err != nil {
			panic(err)
		}

		// .env file does not handle private keys gracefully
		// probably would be better saved to a file than in .env. Oh well.
		invalidKey := os.Getenv("GOOGLE_SHEETS_PRIVATE_KEY")
		fixedKey := strings.ReplaceAll(invalidKey, "__n__", "\n")

		gs := NewClient(Secrets{
			Email:      os.Getenv("GOOGLE_SHEETS_EMAIL"),
			PrivateKey: fixedKey,
		})

		toAdd := []FoundTrackRow{
			{
				Title:   "song 9",
				Artist:  "artist 9",
				AddedAt: "2024-04-16T00:00:00.000Z",
			},
			{
				Title:   "song 2",
				Artist:  "artist 2",
				AddedAt: "2025-04-16T00:00:00.000Z",
			},
		}

		gs.AddFoundTracks(toAdd)
	})

	t.Run("can update 4 source and info columns", func(t *testing.T) {
		t.Skip("skip test calling real google sheets api")

		err := godotenv.Load("../../.env")
		if err != nil {
			panic(err)
		}

		// .env file does not handle private keys gracefully
		// probably would be better saved to a file than in .env. Oh well.
		invalidKey := os.Getenv("GOOGLE_SHEETS_PRIVATE_KEY")
		fixedKey := strings.ReplaceAll(invalidKey, "__n__", "\n")

		gs := NewClient(Secrets{
			Email:      os.Getenv("GOOGLE_SHEETS_EMAIL"),
			PrivateKey: fixedKey,
		})

		values := make([][]interface{}, 4)
		v1 := make([]interface{}, 2)
		v1[0] = "Spotify"
		v1[1] = "Spotify Track information!"
		values[0] = v1

		v2 := make([]interface{}, 2)
		v2[0] = ""
		v2[1] = ""
		values[1] = v2

		v3 := make([]interface{}, 2)
		v3[0] = "GoogleSearch"
		v3[1] = "Found this one from Google"
		values[2] = v3

		v4 := make([]interface{}, 2)
		v4[0] = "GoogleSearch"
		v4[1] = "Found this one from Google"
		values[3] = v4

		gs.updateValues(TESTTrackSheet, "C2:D", values)
	})

	// in case I wanna not update all cells... but maybe I can still update all cells
	t.Run("can update 4 source and info columns: dynamic notation", func(t *testing.T) {
		t.Skip("skip test calling real google sheets api")

		err := godotenv.Load("../../.env")
		if err != nil {
			panic(err)
		}

		// .env file does not handle private keys gracefully
		// probably would be better saved to a file than in .env. Oh well.
		invalidKey := os.Getenv("GOOGLE_SHEETS_PRIVATE_KEY")
		fixedKey := strings.ReplaceAll(invalidKey, "__n__", "\n")

		gs := NewClient(Secrets{
			Email:      os.Getenv("GOOGLE_SHEETS_EMAIL"),
			PrivateKey: fixedKey,
		})

		values := make([][]interface{}, 4)
		v1 := make([]interface{}, 2)
		v1[0] = "Spotify"
		v1[1] = "Spotify Track information!"
		values[0] = v1

		v2 := make([]interface{}, 2)
		v2[0] = ""
		v2[1] = ""
		values[1] = v2

		v3 := make([]interface{}, 2)
		v3[0] = "GoogleSearch"
		v3[1] = "Found this one from Google"
		values[2] = v3

		v4 := make([]interface{}, 2)
		v4[0] = "GoogleSearch"
		v4[1] = "Found this one from Google"
		values[3] = v4

		notation := fmt.Sprintf("C%d:D", 2+3)

		gs.updateValues(TESTTrackSheet, notation, values)
	})
}
