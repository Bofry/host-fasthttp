package failure

import "time"

const (
	INVALID_ARGUMENT  = "INVALID_ARGUMENT"
	INVALID_OPERATION = "INVALID_OPERATION"
	UNKNOWN_ERROR     = "UNKNOWN_ERROR"
	NOP               = "NOP"
	NO_CONTENT        = "NO_CONTENT"
)

func IsKnownErrorCode(message string) bool {
	if size := len(message); (size == 0) || (size > 64) {
		return false
	}

	for i := 0; i < len(message); i++ {
		ch := message[i]

		if ch == '_' ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') {
			continue
		} else {
			return false
		}
	}
	return true
}

func ThrowFailure(err error) {
	var failure *Failure
	if IsKnownErrorCode(err.Error()) {
		failure = &Failure{
			Message:   err.Error(),
			Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
		}
	} else {
		failure = &Failure{
			Message:     INVALID_OPERATION,
			Description: err.Error(),
			Timestamp:   time.Now().UnixNano() / int64(time.Millisecond),
		}
	}
	panic(failure)
}

func ThrowFailureMessage(message string, reason string) {
	var failure *Failure
	if IsKnownErrorCode(message) {
		failure = &Failure{
			Message:     message,
			Description: reason,
			Timestamp:   time.Now().UnixNano() / int64(time.Millisecond),
		}
	} else {
		var desc string = message
		if len(reason) > 0 {
			desc += " " + reason
		}
		failure = &Failure{
			Message:     INVALID_OPERATION,
			Description: desc,
			Timestamp:   time.Now().UnixNano() / int64(time.Millisecond),
		}
	}
	panic(failure)
}
