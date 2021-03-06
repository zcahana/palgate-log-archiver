package googlesheets

import (
	"context"
	"fmt"
	"sort"

	"github.com/zcahana/palgate-log-archiver/sink"
	palgate "github.com/zcahana/palgate-sdk"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const (
	spreadsheetId = "1-wN5Xc59Vx0tp2zGitTF6gBNQOR2qO852qLJ82GW2lM"
	sheetName     = "Sheet1"
)

type sheetsSink struct {
	service *sheets.Service
}

func NewSink() (sink.Sink, error) {
	service, err := initSheetsService()
	if err != nil {
		return nil, fmt.Errorf("error initializing Google Sheets service: %v", err)
	}
	s := &sheetsSink{
		service: service,
	}

	return s, nil
}

func initSheetsService() (*sheets.Service, error) {
	config, err := getGoogleConfig()
	if err != nil {
		return nil, err
	}

	client := getClient(config)
	service, err := sheets.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Google Sheets client: %v", err)
	}

	return service, nil
}

func (s *sheetsSink) Receive(records []palgate.LogRecord) (int, error) {
	rows := rowsFromRecords(records)

	topRow, err := s.readTopRow()
	if err != nil {
		return 0, err
	}

	rows = s.selectNewerRows(rows, topRow)

	err = s.writeTopRows(rows)
	if err != nil {
		return 0, err
	}

	return len(rows), nil
}

func (s *sheetsSink) readTopRow() (Row, error) {
	readRange := fmt.Sprintf("%s!A2:H2", sheetName)
	resp, err := s.service.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 {
		return emptyRow, nil
	}
	if len(resp.Values) > 1 {
		return nil, fmt.Errorf("unexpected number of rows returned: %d", len(resp.Values))
	}

	row := Row(resp.Values[0])
	if err := row.validate(); err != nil {
		return nil, fmt.Errorf("error validating top row: %v", err)
	}

	row = row.defaultize()
	return row, nil
}

func (s *sheetsSink) selectNewerRows(rows []Row, pivot Row) []Row {
	// Newer to older
	sort.Sort(sortableRows(rows))

	if pivot.isEmpty() {
		return rows
	}

	newerRows := make([]Row, 0, len(rows))
	for _, row := range rows {
		if row.isAfter(pivot) {
			newerRows = append(newerRows, row)
		}
	}

	return newerRows
}

func (s *sheetsSink) writeTopRows(rows []Row) error {
	_, err := s.service.Spreadsheets.BatchUpdate(spreadsheetId, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				InsertDimension: &sheets.InsertDimensionRequest{
					InheritFromBefore: false,
					Range: &sheets.DimensionRange{
						Dimension:  "ROWS",
						SheetId:    0,
						StartIndex: 1,
						EndIndex:   int64(1 + len(rows)),
					},
				},
			},
		},
	}).Do()
	if err != nil {
		return err
	}

	writeRange := fmt.Sprintf("%s!A2:H%d", sheetName, 1+len(rows))
	_, err = s.service.Spreadsheets.Values.Update(spreadsheetId, writeRange, &sheets.ValueRange{
		MajorDimension: "ROWS",
		Range:          writeRange,
		Values:         rowsToValues(rows),
	}).ValueInputOption("RAW").Do()
	if err != nil {
		return err
	}

	return nil
}
