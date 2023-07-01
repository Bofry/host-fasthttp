package middleware

import (
	"log"
)

var _ LoggingService = new(CompositeLoggingService)

type CompositeLoggingService struct {
	loggingServices []LoggingService
}

func NewCompositeLoggingService(services ...LoggingService) *CompositeLoggingService {
	return &CompositeLoggingService{
		loggingServices: services,
	}
}

// ConfigureLogger implements LoggingService.
func (s *CompositeLoggingService) ConfigureLogger(l *log.Logger) {
	for _, svc := range s.loggingServices {
		svc.ConfigureLogger(l)
	}
}

// CreateEventLog implements LoggingService.
func (s *CompositeLoggingService) CreateEventLog(ev EventEvidence) EventLog {
	var eventlogs []EventLog
	for _, svc := range s.loggingServices {
		eventlogs = append(eventlogs, svc.CreateEventLog(ev))
	}

	return CompositeEventLog{
		eventLogs: eventlogs,
	}
}
