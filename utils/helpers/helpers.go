package helpers

import (
	// Go Internal Packages
	"fmt"
	"time"

	// External Packages
	"github.com/google/uuid"
)

func GenerateRandomID() string {
	return uuid.New().String()
}

func GetCurrentTime() string {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	istTime := time.Now().In(loc).Format("Jan 02 2006 03:04:05 PM")
	return istTime
}

func GetOrderID(id string) string {
	return fmt.Sprintf("ORDER:%s", id)
}
