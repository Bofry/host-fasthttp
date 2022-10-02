package test

import (
	"log"

	fasthttp "github.com/Bofry/host-fasthttp"
)

var _ fasthttp.LoggingService = new(LoggingService)

type LoggingService struct{}

func (s *LoggingService) CreateEventLog() fasthttp.EventLog {
	log.Println("CreateEventLog()")
	return &EventLog{}
}
