package googlesheets

import (
	"fmt"

	"github.com/zcahana/palgate-sdk"
)

const (
	numOfRowElements = 8
)

var emptyRow = Row{}

// Row abstract a sheet row and provides methods to access its named fields
type Row []interface{}

func (row Row) Date() string {
	return row[0].(string)
}

func (row Row) Time() string {
	return row[1].(string)
}

func (row Row) Type() string {
	return row[2].(string)
}

func (row Row) Status() string {
	return row[3].(string)
}

func (row Row) Serial() string {
	return row[4].(string)
}

func (row Row) UserID() string {
	return row[5].(string)
}

func (row Row) LastName() string {
	return row[6].(string)
}

func (row Row) FirstName() string {
	return row[7].(string)
}

func (row Row) isEmpty() bool {
	for _, value := range row {
		if value == nil {
			continue
		}
		if s, ok := value.(string); ok && s == "" {
			continue
		}
		return false
	}
	return true
}

func (row Row) validate() error {
	if len(row) > numOfRowElements {
		return fmt.Errorf("invalid number of row elements: %d", len(row))
	}

	for i, value := range row {
		if value == nil {
			continue
		}

		_, ok := value.(string)
		if !ok {
			return fmt.Errorf("invalid value at row element %d: %s (type %T)", i, value, value)
		}
	}

	return nil
}

func (row Row) defaultize() Row {
	const defaultValue = ""

	newRow := make(Row, numOfRowElements)
	for i := 0; i < numOfRowElements; i++ {
		if i >= len(row) || row[i] == nil {
			newRow[i] = defaultValue
		} else {
			newRow[i] = row[i]
		}
	}

	return newRow
}

func (row Row) isAfter(otherRow Row) bool {
	// Date
	if row.Date() > otherRow.Date() {
		return true
	}
	if row.Date() < otherRow.Date() {
		return false
	}

	// Time
	if row.Time() > otherRow.Time() {
		return true
	}
	if row.Time() < otherRow.Time() {
		return false
	}

	// UserID
	if row.UserID() > otherRow.UserID() {
		return true
	}
	if row.UserID() < otherRow.UserID() {
		return false
	}

	// Status
	if row.Status() > otherRow.Status() {
		return true
	}
	if row.Status() < otherRow.Status() {
		return false
	}

	// Type
	if row.Type() > otherRow.Type() {
		return true
	}
	if row.Type() < otherRow.Type() {
		return false
	}

	return false
}

func rowFromRecord(record *palgate.LogRecord) Row {
	return Row{
		record.Date(),
		record.Time(),
		record.Type.String(),
		record.OperationStatus.String(),
		record.SerialNumber,
		record.UserID,
		record.LastName,
		record.FirstName,
	}
}

func rowsFromRecords(records []palgate.LogRecord) []Row {
	rows := make([]Row, len(records))
	for i, record := range records {
		rows[i] = rowFromRecord(&record)
	}
	return rows
}

func rowsToValues(rows []Row) [][]interface{} {
	values := make([][]interface{}, len(rows))
	for i, row := range rows {
		values[i] = row
	}
	return values
}

type sortableRows []Row

func (rows sortableRows) Len() int {
	return len(rows)
}

func (rows sortableRows) Less(i int, j int) bool {
	return rows[i].isAfter(rows[j])
}

func (rows sortableRows) Swap(i int, j int) {
	rows[i], rows[j] = rows[j], rows[i]
}
