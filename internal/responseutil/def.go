package responseutil

const (
	// response flag
	SUCCESS ResponseFlag = iota
	FAILURE

	UNKNOWN ResponseFlag = -1
)

type (
	ResponseFlag int

	ResponseState interface {
		Flag() ResponseFlag
		StatusCode() int
	}
)
