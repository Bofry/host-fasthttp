package middleware

import (
	"log"
)

var _ LoggingService = NoopLoggingService{}

type NoopLoggingService struct{}

// ConfigureLogger implements LoggingService.
func (NoopLoggingService) ConfigureLogger(*log.Logger) {}

// CreateEventLog implements LoggingService.
func (NoopLoggingService) CreateEventLog(EventEvidence) EventLog {
	return NoopEventLogSingleton
}
