package helper

import (
	"backend/models"

	"time"
)

func CalculationWorkHour(chekin, checkout time.Time) float64 {
	var result float64
	duration := checkout.Sub(chekin)

	minute := int(duration.Minutes()) % 60
	hour := int(duration.Hours())

	if hour < 1 {
		result = float64(minute) / 100
	} else {
		result = float64(hour) + float64(minute)/100
	}

	return result
}

func CalculationWorkMonthly(attedances []models.Attedance) float64 {
	var totalDuration time.Duration

	for _, attedance := range attedances {
		totalDuration += attedance.Chekout.Sub(*attedance.Chekin)
	}

	minute := int(totalDuration.Minutes()) % 60
	hour := int(totalDuration.Hours())
	result := float64(hour) + float64(minute)/100
	return result
}
