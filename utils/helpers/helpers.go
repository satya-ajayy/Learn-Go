package helpers

import (
	"github.com/google/uuid"
	"time"
)

func GenerateRandomID() string {
	return uuid.New().String()
}

func GetCurrentTime() time.Time {
	return time.Now()
}

func GetCurrentTimeString() string {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	istTime := time.Now().In(loc).Format("Jan 02 2006 03:04:05 PM")
	return istTime
}
