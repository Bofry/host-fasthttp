package internal

import "reflect"

type RequestResourceProcessService struct {
	modules []RequestResourceProcessModule
}

func NewRequestHandlerReprocessService() *RequestResourceProcessService {
	return &RequestResourceProcessService{}
}

func (s *RequestResourceProcessService) Register(module RequestResourceProcessModule) {
	s.modules = append(s.modules, module)
}

func (s *RequestResourceProcessService) Process(requestManager interface{}) {
	var (
		rvManager reflect.Value = reflect.ValueOf(requestManager)
	)
	if rvManager.Kind() != reflect.Pointer || rvManager.IsNil() {
		return
	}

	rvManager = reflect.Indirect(rvManager)
	numOfHandles := rvManager.NumField()
	for i := 0; i < numOfHandles; i++ {
		rvHandler := rvManager.Field(i)
		if rvHandler.Kind() != reflect.Pointer || rvHandler.IsNil() {
			continue
		}

		rvHandler = reflect.Indirect(rvHandler)
		if rvHandler.Kind() == reflect.Struct {
			for _, v := range s.modules {
				v.ProcessRequestResource(rvHandler)
			}
		}
	}
}
