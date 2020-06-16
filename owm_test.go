package owm

import (
	"testing"
	"os"
)

func TestGetWeather(t *testing.T) {
	w, err := GetWeather(52.41467, 10.74063, os.Getenv("OWM_API_KEY"))
	if err != nil {
		t.Error(err)
		return
	}

	if len(w.DailyWeather) < 1 {
		t.Error("No daily weather received")
	}

	if len(w.HourlyWeather) < 1 {
		t.Error("No hourly weather received")
	}
}