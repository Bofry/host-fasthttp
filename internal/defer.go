package internal

import (
	"sync/atomic"
)

type Defer struct {
	recover *Recover
	catch   func(err interface{})
	finally []func(err interface{})
}

func (d *Defer) Do(do func(f Finalizer)) {
	if d != nil {
		defer func(f Finalizer) {
			var err interface{} = nil
			if d.recover != nil {
				if atomic.LoadUint32(&d.recover.done) == 0 {
					atomic.StoreUint32(&d.recover.done, 1)
					d.recover.err = recover()
				}
				err = d.recover.err
			}
			d.catch(err)
			f.run(err)
		}(Finalizer{d})
	}
	do(Finalizer{d})
}

type Finalizer struct {
	deferService *Defer
}

func (f Finalizer) Add(actions ...func(err interface{})) {
	f.deferService.finally = append(f.deferService.finally, actions...)
}

func (f Finalizer) run(err interface{}) {
	for _, h := range f.deferService.finally {
		h(err)
	}
}

type Recover struct {
	err  interface{}
	done uint32
}

func (s *Recover) Defer(catch func(err interface{})) *Defer {
	return &Defer{
		recover: s,
		catch:   catch,
	}
}
