package internal

import "sync/atomic"

type DeferService struct {
	recover *RecoverService
	finally func(err interface{})
}

func (s *DeferService) Do(f func()) {
	if s != nil {
		defer func() {
			if s.recover != nil {
				if atomic.LoadUint32(&s.recover.done) == 0 {
					atomic.StoreUint32(&s.recover.done, 1)
					s.recover.err = recover()
				}
				s.finally(s.recover.err)
			} else {
				s.finally(nil)
			}
		}()
	}
	f()
}
