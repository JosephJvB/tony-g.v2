package googlesheets

import (
	"strconv"
)

type TonyVideoRow struct {
	Id          string
	Title       string
	PublishedAt string
	TotalTracks int
	AddedAt     string
}
type FoundTrackRow struct {
	Title                  string
	Artist                 string
	FoundVideoTitle        string
	FoundChannelTitle      string
	FoundVideoId           string
	Link                   string
	Confidence             string
	ReviewVideoId          string
	ReviewVideoPublishDate string
	AddedAt                string
	Playlist               string
}

func RowToTonyVideo(row []interface{}) TonyVideoRow {
	ttStr := row[3].(string)
	tt, err := strconv.Atoi(ttStr)
	if err != nil {
		tt = -1
	}

	return TonyVideoRow{
		Id:          row[0].(string),
		Title:       row[1].(string),
		PublishedAt: row[2].(string),
		TotalTracks: tt,
		AddedAt:     row[4].(string),
	}
}

func TonyVideoToRow(video TonyVideoRow) []interface{} {
	r := make([]interface{}, 5)
	r[0] = video.Id
	r[1] = video.Title
	r[2] = video.PublishedAt
	r[3] = video.TotalTracks
	r[4] = video.AddedAt

	return r
}

func RowToFoundTrack(row []interface{}) FoundTrackRow {
	return FoundTrackRow{
		Title:                  row[0].(string),
		Artist:                 row[1].(string),
		FoundVideoTitle:        row[2].(string),
		FoundChannelTitle:      row[3].(string),
		FoundVideoId:           row[4].(string),
		Link:                   row[5].(string),
		Confidence:             row[6].(string),
		ReviewVideoId:          row[7].(string),
		ReviewVideoPublishDate: row[8].(string),
		AddedAt:                row[9].(string),
		Playlist:               row[10].(string),
	}
}

func FoundTrackToRow(track FoundTrackRow) []interface{} {
	r := make([]interface{}, 11)
	r[0] = track.Title
	r[1] = track.Artist
	r[2] = track.FoundVideoTitle
	r[3] = track.FoundChannelTitle
	r[4] = track.FoundVideoId
	r[5] = track.Link
	r[6] = track.Confidence
	r[7] = track.ReviewVideoId
	r[8] = track.ReviewVideoPublishDate
	r[9] = track.AddedAt
	r[10] = track.Playlist

	return r
}
