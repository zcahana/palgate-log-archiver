package palgate

import (
	"fmt"
	"time"
)

const (
	ResponseStatusSuccess = "ok"
	ResponseStatusFailed  = "failed"
)

type GetLogResponse struct {
	Records []LogRecord `json:"log"`
	Error   string      `json:"err"`
	Message string      `json:"msg"`
	Status  string      `json:"status"`
}

type LogRecord struct {
	UserID          string          `json:"userId"`
	OperationStatus OperationStatus `json:"operation"`
	Timestamp       int             `json:"time"`
	FirstName       string          `json:"firstname,omitempty"`
	LastName        string          `json:"lastname,omitempty"`
	Image           bool            `json:"image"`
	Reason          int             `json:"reason"`
	Type            OperationType   `json:"type"`
	SerialNumber    string          `json:"sn"`
}

type OperationStatus string

const (
	OperationStatusSuccess   = "sr1"
	OperationStatusBadSignal = "sr13"
	OperationStatusUndefined = "srundefined"
)

func (s OperationStatus) String() string {
	switch s {
	case OperationStatusSuccess:
		return "Success"
	case OperationStatusBadSignal:
		return "Bad Signal"
	case OperationStatusUndefined:
		return "Undefined"
	default:
		return "Unknown"
	}
}

type OperationType int

const (
	OperationDial          OperationType = 1
	OperationRemoteControl OperationType = 2
	OperationApplication   OperationType = 100
)

func (o OperationType) String() string {
	switch o {
	case OperationDial:
		return "Dial"
	case OperationRemoteControl:
		return "Remote Control"
	case OperationApplication:
		return "Application"
	default:
		return "Unknown"
	}
}

func (record *LogRecord) Name() string {
	return fmt.Sprintf("%s %s", record.FirstName, record.LastName)
}

func (record *LogRecord) Date() string {
	const format = "2/1/2006"
	t := time.Unix(int64(record.Timestamp), 0)
	return t.Local().Format(format)
}

func (record *LogRecord) Time() string {
	const format = "15:04:05"
	t := time.Unix(int64(record.Timestamp), 0)
	return t.Local().Format(format)
}

type UserRequest struct {
	ID           string `json:"id"`
	FirstName    string `json:"firstname"`
	LastName     string `json:"lastname"`
	Admin        bool   `json:"admin"`
	DialToOpen   bool   `json:"dialToOpen"`
	Output1      bool   `json:"output1"`
	Output2      bool   `json:"output2"`
	Output1Latch bool   `json:"output1Latch"`
	Output2Latch bool   `json:"output2Latch"`
}

type GetUsersResponse struct {
	Error   string `json:"err"`
	Message string `json:"msg"`
	Status  string `json:"status"`
	Start   int    `json:"start"`
	Length  int    `json:"len"`
	Count   int    `json:"count"`
	Users   []User `json:"users"`
}

type User struct {
	ID                  string `json:"id"`
	FirstName           string `json:"firstname"`
	LastName            string `json:"lastname"`
	Admin               bool   `json:"admin"`
	DialToOpen          bool   `json:"dialToOpen"`
	AppInstalled        bool   `json:"appInstalled"`
	Image               bool   `json:"image"`
	Output1             bool   `json:"output1"`
	Output2             bool   `json:"output2,omitempty"`
	Output1Latch        bool   `json:"output1Latch"`
	Output2Latch        bool   `json:"output2Latch"`
	Output1LatchMaxTime int    `json:"output1LatchMaxTime"`
	Output2LatchMaxTime int    `json:"output2LatchMaxTime"`
	Output1Color        string `json:"output1Color"`
	Output2Color        string `json:"output2Color,omitempty"`
	Output1Icon         string `json:"output1Icon"`
	Output2Icon         string `json:"output2Icon,omitempty"`
}
