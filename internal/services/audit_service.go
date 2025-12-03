package services

import (
	"github.com/rs/zerolog/log"
	"time"
)

type AuditLog struct {
	Action    string    `json:"action"`
	UserID    int       `json:"user_id,omitempty"`
	Data      string    `json:"data,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Error     string    `json:"error,omitempty"`
}

type AuditService struct {
	logChan chan AuditLog
}

func NewAuditService() *AuditService {
	service := &AuditService{
		logChan: make(chan AuditLog, 1000),
	}
	go service.processLogs()
	return service
}

func (s *AuditService) LogAsync(action string, userID int, data string) {
	go s.log(action, userID, data, "")
}

func (s *AuditService) LogErrorAsync(action string, userID int, err error) {
	go s.log(action, userID, "", err.Error())
}

func (s *AuditService) log(action string, userID int, data, err string) {
	select {
	case s.logChan <- AuditLog{
		Action:    action,
		UserID:    userID,
		Data:      data,
		Timestamp: time.Now(),
		Error:     err,
	}:
	default:
		log.Warn().Msg("Audit log channel is full")
	}
}

func (s *AuditService) processLogs() {
	for logEntry := range s.logChan {
		log.Info().
			Str("action", logEntry.Action).
			Int("user_id", logEntry.UserID).
			Str("data", logEntry.Data).
			Str("error", logEntry.Error).
			Time("timestamp", logEntry.Timestamp).
			Msg("Audit log")
	}
}
