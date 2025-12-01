package googlesheets

import (
	"context"
	"log"

	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const SpreadsheetId = "1PFjWpTSX5iZH-0B1yUHSuqO9ixGW3LnBLFiv52bYcIg"

type SheetConfig struct {
	Name        string
	Id          int
	AllRowRange string
}

var TonysVideoSheet = SheetConfig{
	Name:        "Tony's Videos",
	Id:          571234434,
	AllRowRange: "A2:F",
}
var FoundTrackSheet = SheetConfig{
	Name:        "Found Tracks",
	Id:          1461864733,
	AllRowRange: "A2:I",
}
var TESTTrackSheet = SheetConfig{
	Name:        "TEST",
	Id:          407333682,
	AllRowRange: "A2:I",
}

type GoogleSheetsClient struct {
	sheetsService *sheets.Service
}

type Secrets struct {
	Email      string
	PrivateKey string
}

// https://gist.github.com/karayel/1b915b61d3cf307ca23b14313848f3c4
func NewClient(secrets Secrets) GoogleSheetsClient {
	conf := &jwt.Config{
		Email:      secrets.Email,
		PrivateKey: []byte(secrets.PrivateKey),
		TokenURL:   "https://oauth2.googleapis.com/token",
		Scopes: []string{
			"https://www.googleapis.com/auth/spreadsheets",
		},
	}

	client := conf.Client(context.Background())

	sheetsService, err := sheets.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		panic(err)
	}

	return GoogleSheetsClient{
		sheetsService: sheetsService,
	}
}

func (gs *GoogleSheetsClient) GetTonysVideos() []TonyVideoRow {
	rows := gs.getRows(TonysVideoSheet)

	videos := []TonyVideoRow{}
	for _, row := range rows {
		r := RowToTonyVideo(row)

		videos = append(videos, r)
	}

	return videos
}

func (gs *GoogleSheetsClient) AddTonysVideos(nextRows []TonyVideoRow) {
	// sheets.ValueRange.Values needs interfaces
	rows := make([][]interface{}, len(nextRows))
	for _, t := range nextRows {
		r := TonyVideoToRow(t)

		rows = append(rows, r)
	}

	gs.appendRows(TonysVideoSheet, rows)
}

func (gs *GoogleSheetsClient) GetFoundTracks() []FoundTrackRow {
	rows := gs.getRows(FoundTrackSheet)

	tracks := []FoundTrackRow{}
	for _, row := range rows {
		r := RowToFoundTrack(row)

		tracks = append(tracks, r)
	}

	return tracks
}

func (gs *GoogleSheetsClient) AddFoundTracks(nextRows []FoundTrackRow) {
	// sheets.ValueRange.Values needs interfaces
	rows := make([][]interface{}, len(nextRows))
	for _, t := range nextRows {
		r := FoundTrackToRow(t)

		rows = append(rows, r)
	}

	gs.appendRows(FoundTrackSheet, rows)
}

func (gs *GoogleSheetsClient) getRows(cfg SheetConfig) [][]interface{} {
	sheetRange := cfg.Name + "!" + cfg.AllRowRange

	resp, err := gs.sheetsService.Spreadsheets.Values.
		Get(SpreadsheetId, sheetRange).
		Do()
	if err != nil {
		// I think these errors were due to Office Wi-Fi dropping
		// happened twice in a row! And then stopped
		// something like this
		// googleapi: Error 500: Internal error encountered., backendError
		// os.WriteFile("./data/googlesheets-error.txt", []byte(err.Error()), 0666)
		log.Fatal(err)
	}

	return resp.Values
}

func (gs *GoogleSheetsClient) appendRows(cfg SheetConfig, rows [][]interface{}) {
	// set next rows
	valueRange := sheets.ValueRange{
		MajorDimension: "ROWS",
		Values:         rows,
	}
	// is this range gonna append rows the way I want?
	rowRange := cfg.Name + "!" + cfg.AllRowRange
	req := gs.sheetsService.Spreadsheets.Values.Append(SpreadsheetId, rowRange, &valueRange)
	// is this the only way to add these params?
	req.ValueInputOption("RAW")
	// other option is "OVERWRITE"
	// but that only overwrites if there's empty cells, not what I expected
	// I guess it's "Append" method after all
	req.InsertDataOption("INSERT_ROWS")

	req.Do()
}

func (gs *GoogleSheetsClient) updateValues(cfg SheetConfig, cellRange string, values [][]interface{}) {
	valueRange := sheets.ValueRange{
		MajorDimension: "ROWS",
		Values:         values,
	}

	updateRange := cfg.Name + "!" + cellRange

	req := gs.sheetsService.Spreadsheets.Values.Update(
		SpreadsheetId,
		updateRange,
		&valueRange,
	)
	req.ValueInputOption("RAW")
	req.Do()
}
