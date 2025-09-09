package failure

import (
	"encoding/json"
	"fmt"
)

var _ error = new(Failure)

type Failure struct {
	Message     string          `json:"message"`
	Description string          `json:"description,omitempty"`
	Reason      json.RawMessage `json:"reason,omitempty"`
	Timestamp   int64           `json:"timestamp"`
	TraceId     string          `json:"_trace_id,omitempty"`
	Err         error           `json:"-"`
}

func (f *Failure) Error() string {
	if len(f.Description) > 0 {
		return fmt.Sprintf("%s - %s", f.Message, f.Description)
	}
	return f.Message
}

func (f *Failure) Unwrap() error { return f.Err }

func (f *Failure) MarshalJSON() ([]byte, error) {
	var reason *string = nil
	if len(f.Description) > 0 {
		reason = &f.Description
	}

	type Alias Failure
	return json.Marshal(&struct {
		Description *string `json:"description"`
		*Alias
	}{
		Description: reason,
		Alias:       (*Alias)(f),
	})
}
