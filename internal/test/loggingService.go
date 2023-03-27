package test

import (
	"log"

	fasthttp "github.com/Bofry/host-fasthttp"
)

var _ fasthttp.LoggingService = new(LoggingService)

type LoggingService struct {
	logger *log.Logger
}

func (s *LoggingService) CreateEventLog(ev fasthttp.EventEvidence) fasthttp.EventLog {
	s.logger.Println("CreateEventLog()")
	return EventLog{
		logger:   s.logger,
		evidence: ev,
	}
}

func (s *LoggingService) ConfigureLogger(l *log.Logger) {
	s.logger = l
}
