package helper

import (
	"time"
)

func MonthsSinceJoined(createdAt time.Time) int {
	now := time.Now()
	years := now.Year() - createdAt.Year()
	months := int(now.Month()) - int(createdAt.Month())
	if months < 0 {
		years--
		months += 12
	}
	cal := years*12 + months

	if cal == 0 {
		return 1
	} else {
		return cal

	}
}
