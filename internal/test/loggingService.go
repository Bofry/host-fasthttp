package test

import (
	"fmt"

	fasthttp "github.com/Bofry/host-fasthttp"
)

var _ fasthttp.LoggingService = new(LoggingService)

type LoggingService struct{}

func (s *LoggingService) CreateEventLog() fasthttp.EventLog {
	fmt.Println("CreateEventLog()")
	return &EventLog{}
}
