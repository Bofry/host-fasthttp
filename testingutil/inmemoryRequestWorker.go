package testingutil

type inmemoryRequestWorker struct {
	ErrorChannel chan interface{}
}

func (w *inmemoryRequestWorker) NoticeError(err interface{}) {
	if err != nil {
		if w.ErrorChannel != nil {
			w.ErrorChannel <- err
		}
	}
}
