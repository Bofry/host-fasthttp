package internal

type RecoverService struct {
	err  interface{}
	done uint32
}

func (s *RecoverService) Defer(catch func(err interface{})) *DeferService {
	return &DeferService{
		recover: s,
		catch:   catch,
	}
}
