package googlesheets

import (
	"log"
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
	FoundTrackInfo         string
	TrackVideoId           string
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
	log.Printf("%v", row)
	return FoundTrackRow{
		Title:                  row[0].(string),
		Artist:                 row[1].(string),
		FoundTrackInfo:         row[2].(string),
		TrackVideoId:           row[3].(string),
		Link:                   row[4].(string),
		Confidence:             row[5].(string),
		ReviewVideoId:          row[6].(string),
		ReviewVideoPublishDate: row[7].(string),
		AddedAt:                row[8].(string),
		Playlist:               row[9].(string),
	}
}

func FoundTrackToRow(track FoundTrackRow) []interface{} {
	r := make([]interface{}, 10)
	r[0] = track.Title
	r[1] = track.Artist
	r[2] = track.FoundTrackInfo
	r[3] = track.TrackVideoId
	r[4] = track.Link
	r[5] = track.Confidence
	r[6] = track.ReviewVideoId
	r[7] = track.ReviewVideoPublishDate
	r[8] = track.AddedAt
	r[9] = track.Playlist

	return r
}
