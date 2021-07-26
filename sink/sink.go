package sink

import (
	"github.com/zcahana/palgate-sdk"
)

type Sink interface {
	Receive(records []palgate.LogRecord) (int, error)
}
