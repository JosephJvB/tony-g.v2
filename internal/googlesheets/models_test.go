package googlesheets

import (
	"testing"
)

func TestGoogleSheetsModels(t *testing.T) {
	t.Run("Row to Tony Video", func(t *testing.T) {
		r := make([]interface{}, 5)
		r[0] = "id-123"
		r[1] = "weekly tracko roundo"
		r[2] = "2025-01-01"
		r[3] = "10"
		r[4] = "ages ago"

		v := RowToTonyVideo(r)

		if v.Id != "id-123" {
			t.Errorf("expected id to be id-123. Got %s", v.Id)
		}
		if v.Title != "weekly tracko roundo" {
			t.Errorf("expected title to be weekly tracko roundo. Got %s", v.Title)
		}
		if v.PublishedAt != "2025-01-01" {
			t.Errorf("expected publishedAt to be 2025-01-01. Got %s", v.PublishedAt)
		}
		if v.TotalTracks != 10 {
			t.Errorf("expected totalTracks to be 10. Got %d", v.TotalTracks)
		}
		if v.AddedAt != "ages ago" {
			t.Errorf("expected addedAt to be ages ago. Got %s", v.AddedAt)
		}
	})

	t.Run("Youtube Video To Row", func(t *testing.T) {
		v := TonyVideoRow{
			Id:          "id-123",
			Title:       "weekly tracko roundo",
			PublishedAt: "2025-01-01",
			TotalTracks: 10,
			AddedAt:     "ages ago",
		}

		row := TonyVideoToRow(v)

		if len(row) != 5 {
			t.Errorf("expected row to have 5 elements. Got %d", len(row))
		}

		if row[0] != "id-123" {
			t.Errorf("expected row[0] to be id-123. Got %s", row[0])
		}
		if row[1] != "weekly tracko roundo" {
			t.Errorf("expected row[1] to be weekly tracko roundo. Got %s", row[1])
		}
		if row[2] != "2025-01-01" {
			t.Errorf("expected row[2] to be 2025-01-01. Got %s", row[2])
		}
		if row[3] != 10 {
			t.Errorf("expected row[3] to be 10. Got %d", row[3])
		}
		if row[4] != "ages ago" {
			t.Errorf("expected row[5] to be ages ago. Got %s", row[5])
		}
	})

	t.Run("Youtube Track to Row", func(t *testing.T) {
		yt := FoundTrackRow{
			Title:                  "little things",
			Artist:                 "adrianne",
			FoundVideoTitle:        "adrianne little things",
			FoundChannelTitle:      "adriannes channel",
			FoundVideoId:           "123",
			Link:                   "https://www.youtube.com/watch?v=123",
			Confidence:             "1.2",
			ReviewVideoId:          "456",
			ReviewVideoPublishDate: "recently",
			AddedAt:                "ages ago",
			Playlist:               "2024",
		}

		row := FoundTrackToRow(yt)

		if len(row) != 11 {
			t.Errorf("expected row to have 11 elements. Got %d", len(row))
		}

		if row[0] != "little things" {
			t.Errorf("expected row[0] to be little things. Got %s", row[0])
		}
		if row[1] != "adrianne" {
			t.Errorf("expected row[1] to be adrianne. Got %s", row[1])
		}
		if row[2] != "adrianne little things" {
			t.Errorf("expected row[2] to be adrianne little things, Got %s", row[2])
		}
		if row[3] != "adriannes channel" {
			t.Errorf("expected row[3] to be adriannes channel. Got %s", row[3])
		}
		if row[4] != "123" {
			t.Errorf("expected row[4] to be 123. Got %s", row[4])
		}
		if row[5] != "https://www.youtube.com/watch?v=123" {
			t.Errorf("expected row[5] to be https://www.youtube.com/watch?v=123. Got %s", row[5])
		}
		if row[6] != "1.2" {
			t.Errorf("Expected row[6] to be 1.2. Got %s", row[6])
		}
		if row[7] != "456" {
			t.Errorf("expected row[7] to be 456. Got %s", row[7])
		}
		if row[8] != "recently" {
			t.Errorf("expected row[8] to be recently. Got %s", row[8])
		}
		if row[9] != "ages ago" {
			t.Errorf("expected row[9] to be ages ago. Got %s", row[9])
		}
		if row[10] != "2024" {
			t.Errorf("expected row[10] to be 2024. Got %s", row[10])
		}
	})

	t.Run("Row to Youtube Track", func(t *testing.T) {
		r := make([]interface{}, 11)
		r[0] = "little things"
		r[1] = "adrianne"
		r[2] = "adrianne little things"
		r[3] = "adriannes channel"
		r[4] = "123"
		r[5] = "https://www.youtube.com/watch?v=123"
		r[6] = "1.2"
		r[7] = "456"
		r[8] = "recently"
		r[9] = "ages ago"
		r[10] = "2024"

		yt := RowToFoundTrack(r)

		if yt.Title != "little things" {
			t.Errorf("expected title to be little things. Got %s", yt.Title)
		}
		if yt.Artist != "adrianne" {
			t.Errorf("expected artist to be adrianne. Got %s", yt.Artist)
		}
		if yt.FoundVideoTitle != "adrianne little things" {
			t.Errorf("expected foundVideoTitle to be adrianne little things. Got %s", yt.FoundVideoTitle)
		}
		if yt.FoundChannelTitle != "adriannes channel" {
			t.Errorf("expected foundChannelTitle to be adriannes channel. Got %s", yt.FoundChannelTitle)
		}
		if yt.FoundVideoId != "123" {
			t.Errorf("expected TrackVideoId to be 123. Got %s", yt.FoundVideoId)
		}
		if yt.Link != "https://www.youtube.com/watch?v=123" {
			t.Errorf("expected link to be https://www.youtube.com/watch?v=123. Got %s", yt.Link)
		}
		if yt.Confidence != "1.2" {
			t.Errorf("expected ReviewVideoId to be 1.2. Got %s", yt.Confidence)
		}
		if yt.ReviewVideoId != "456" {
			t.Errorf("expected ReviewVideoId to be 456. Got %s", yt.ReviewVideoId)
		}
		if yt.ReviewVideoPublishDate != "recently" {
			t.Errorf("expected ReviewVideoPublishDate to be recently. Got %s", yt.ReviewVideoPublishDate)
		}
		if yt.AddedAt != "ages ago" {
			t.Errorf("expected addedAt to be ages ago. Got %s", yt.AddedAt)
		}
		if yt.Playlist != "2024" {
			t.Errorf("expected playlist to be 2024. Got %s", yt.Playlist)
		}
	})
}
