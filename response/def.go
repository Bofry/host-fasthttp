package response

const (
	// response flag
	SUCCESS ResponseFlag = iota
	FAILURE

	UNKNOWN ResponseFlag = -1
)

type ResponseFlag int

type Response interface {
	Flag() ResponseFlag
	StatusCode() int
}
