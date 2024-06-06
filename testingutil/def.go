package testingutil

type (
	InmemoryRequestProcessingOption interface {
		apply(*inmemoryRequestWorker) error
	}
)
