package logger

import (
	// Go Internal Packages
	"context"
	"fmt"

	// Local Packages
	models "learn-go/models/logging"
	helpers "learn-go/utils/helpers"
)

type LogsRepository interface {
	InsertLog(ctx context.Context, student models.LogModel) error
}

type LoggerService struct {
	loggerRepository LogsRepository
}

func NewService(loggerRepository LogsRepository) *LoggerService {
	return &LoggerService{loggerRepository: loggerRepository}
}

func (s *LoggerService) InsertLog(ctx context.Context, log models.LogModel) (string, error) {
	log.ID = helpers.GenerateRandomID()
	log.Timestamp = helpers.GetCurrentTimeString()
	err := s.loggerRepository.InsertLog(ctx, log)
	if err != nil {
		return "", fmt.Errorf("failed to insert log due to :: %w", err)
	}
	return log.ID, nil
}
