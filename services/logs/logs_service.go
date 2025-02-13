package logs

import (
	// Go Internal Packages
	"context"
	"fmt"

	// Local Packages
	models "learn-go/models/logs"
	helpers "learn-go/utils/helpers"
)

type LogsRepository interface {
	InsertLog(ctx context.Context, student models.LogModel) error
}

type LogsService struct {
	loggerRepository LogsRepository
}

func NewService(loggerRepository LogsRepository) *LogsService {
	return &LogsService{loggerRepository: loggerRepository}
}

func (s *LogsService) InsertLog(ctx context.Context, log models.LogModel) (string, error) {
	log.ID = helpers.GenerateRandomID()
	log.Timestamp = helpers.GetCurrentTimeString()
	err := s.loggerRepository.InsertLog(ctx, log)
	if err != nil {
		return "", fmt.Errorf("failed to insert log due to :: %w", err)
	}
	return log.ID, nil
}
