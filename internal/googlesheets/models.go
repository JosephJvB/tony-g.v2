package googlesheets

import (
	"strconv"
)

type TonyVideoRow struct {
	Id          string
	Title       string
	PublishedAt string
	TotalTracks int
	FoundTracks int
	AddedAt     string
}
type FoundTrackRow struct {
	Title                  string
	Artist                 string
	FoundTrackInfo         string
	TrackVideoId           string
	Link                   string
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
	ftStr := row[4].(string)
	ft, err := strconv.Atoi(ftStr)
	if err != nil {
		ft = -1
	}

	return TonyVideoRow{
		Id:          row[0].(string),
		Title:       row[1].(string),
		PublishedAt: row[2].(string),
		TotalTracks: tt,
		FoundTracks: ft,
		AddedAt:     row[5].(string),
	}
}

func TonyVideoToRow(video TonyVideoRow) []interface{} {
	r := make([]interface{}, 6)
	r[0] = video.Id
	r[1] = video.Title
	r[2] = video.PublishedAt
	r[3] = video.TotalTracks
	r[4] = video.FoundTracks
	r[5] = video.AddedAt

	return r
}

func RowToFoundTrack(row []interface{}) FoundTrackRow {
	return FoundTrackRow{
		Title:                  row[0].(string),
		Artist:                 row[1].(string),
		FoundTrackInfo:         row[2].(string),
		TrackVideoId:           row[3].(string),
		Link:                   row[4].(string),
		ReviewVideoId:          row[5].(string),
		ReviewVideoPublishDate: row[6].(string),
		AddedAt:                row[7].(string),
		Playlist:               row[8].(string),
	}
}

func FoundTrackToRow(track FoundTrackRow) []interface{} {
	r := make([]interface{}, 9)
	r[0] = track.Title
	r[1] = track.Artist
	r[2] = track.FoundTrackInfo
	r[3] = track.TrackVideoId
	r[4] = track.Link
	r[5] = track.ReviewVideoId
	r[6] = track.ReviewVideoPublishDate
	r[7] = track.AddedAt
	r[8] = track.Playlist

	return r
}
