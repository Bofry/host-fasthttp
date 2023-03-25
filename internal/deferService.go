package internal

import (
	"sync/atomic"
)

type DeferService struct {
	recover *RecoverService
	catch   func(err interface{})
	finally []func(err interface{})
}

func (s *DeferService) Do(f func(f Finalizer)) {
	if s != nil {
		defer func(f Finalizer) {
			var err interface{} = nil
			if s.recover != nil {
				if atomic.LoadUint32(&s.recover.done) == 0 {
					atomic.StoreUint32(&s.recover.done, 1)
					s.recover.err = recover()
				}
				err = s.recover.err
			}
			f.run(err)
			s.catch(err)
		}(Finalizer{s})
	}
	f(Finalizer{s})
}

type Finalizer struct {
	deferService *DeferService
}

func (f Finalizer) Add(actions ...func(err interface{})) {
	f.deferService.finally = append(f.deferService.finally, actions...)
}

func (f Finalizer) run(err interface{}) {
	for _, h := range f.deferService.finally {
		h(err)
	}
}
