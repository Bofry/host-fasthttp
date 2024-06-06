package testingutil

var _ InmemoryRequestProcessingOption = InmemoryRequestProcessingOptionFunc(nil)

type InmemoryRequestProcessingOptionFunc func(*inmemoryRequestWorker) error

func (fn InmemoryRequestProcessingOptionFunc) apply(c *inmemoryRequestWorker) error {
	return fn(c)
}

// --------------------------------------------
func WithErrorChannel(ch chan interface{}) InmemoryRequestProcessingOption {
	return InmemoryRequestProcessingOptionFunc(func(opt *inmemoryRequestWorker) error {
		if ch != nil {
			opt.ErrorChannel = ch
		}
		return nil
	})
}
