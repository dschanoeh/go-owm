package owm

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	baseURL = "http://api.openweathermap.org/data/2.5/onecall?"
)

type WeatherData struct {
	Latitude       float64              `json:"lat"`
	Longitude      float64              `json:"lon"`
	Timezone       string               `json:"timezone"`
	TimezoneOffset uint16               `json:"timezone_offset"`
	Current        CurrentData          `json:"current"`
	HourlyWeather  []HourlyWeatherSlice `json:"hourly"`
	DailyWeather   []DailyWeatherSlice  `json:"daily"`
	ParsedTimeZone *time.Location
}

type CurrentData struct {
	Timestamp     int64                `json:"dt"`
	Sunrise       int64                `json:"sunrise"`
	Sunset        int64                `json:"sunset"`
	Temperature   float64              `json:"temp"`
	FeelsLike     float64              `json:"feels_like"`
	Pressure      float64              `json:"pressure"`
	Humidity      uint8                `json:"humidity"`
	DewPoint      float64              `json:"dew_point"`
	UVI           float64              `json:"uvi"`
	Clouds        uint8                `json:"clouds"`
	Visiblity     uint64               `json:"visibility"`
	WindSpeed     float64              `json:"wind_speed"`
	WindDirection float64              `json:"wind_deg"`
	Weather       []WeatherDescription `json:"weather"`
	Rain          Rain                 `json:"rain"`
	Snow          Snow                 `json:"snow"`
	ParsedTime    *time.Time
}

type WeatherDescription struct {
	ID          uint16 `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type HourlyWeatherSlice struct {
	Timestamp     int64                `json:"dt"`
	Temperature   float64              `json:"temp"`
	FeelsLike     float64              `json:"feels_like"`
	Pressure      float64              `json:"pressure"`
	Humidity      uint8                `json:"humidity"`
	DewPoint      float64              `json:"dew_point"`
	Clouds        uint8                `json:"clouds"`
	WindSpeed     float64              `json:"wind_speed"`
	WindDirection float64              `json:"wind_deg"`
	Weather       []WeatherDescription `json:"weather"`
	Rain          Rain                 `json:"rain"`
	Snow          Snow                 `json:"snow"`
	ParsedTime    *time.Time
}

type DailyWeatherSlice struct {
	Timestamp     int64                `json:"dt"`
	Sunrise       int64                `json:"sunrise"`
	Sunset        int64                `json:"sunset"`
	Temperature   TemperatureStats     `json:"temp"`
	FeelsLike     FeelsLikeStats       `json:"feels_like"`
	Pressure      float64              `json:"pressure"`
	Humidity      uint8                `json:"humidity"`
	DewPoint      float64              `json:"dew_point"`
	WindSpeed     float64              `json:"wind_speed"`
	WindDirection float64              `json:"wind_deg"`
	Weather       []WeatherDescription `json:"weather"`
	Clouds        uint8                `json:"clouds"`
	Rain          float64              `json:"rain"`
	Snow          float64              `json:"snow"`
	UVI           float64              `json:"uvi"`
	ParsedTime    *time.Time
}

type TemperatureStats struct {
	Day     float64 `json:"day"`
	Min     float64 `json:"min"`
	Max     float64 `json:"max"`
	Night   float64 `json:"night"`
	Evening float64 `json:"eve"`
	Morning float64 `json:"morn"`
}

type FeelsLikeStats struct {
	Day     float64 `json:"day"`
	Night   float64 `json:"night"`
	Evening float64 `json:"eve"`
	Morning float64 `json:"morn"`
}

type Rain struct {
	OneHour float64 `json:"1h"`
}

type Snow struct {
	OneHour float64 `json:"1h"`
}

// GetWeather performs an API call for the location given through lat,lon using the APIkey
func GetWeather(lat float64, lon float64, APIKey string) (*WeatherData, error) {
	resp, err := http.Get(fmt.Sprintf("%sappid=%s&lat=%f&lon=%f&units=metric", baseURL, APIKey, lat, lon))

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, errors.New("Received status code 401 - the API key is likely invalid")
		}

		return nil, fmt.Errorf("Received status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	d := WeatherData{}

	if err = json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return nil, err
	}

	parseData(&d)

	return &d, nil
}

func parseData(w *WeatherData) error {
	parsedTimeZone, err := time.LoadLocation(w.Timezone)
	if err != nil {
		return err
	}
	w.ParsedTimeZone = parsedTimeZone

	currentTime := time.Unix(w.Current.Timestamp, 0)
	w.Current.ParsedTime = &currentTime

	for _, data := range w.HourlyWeather {
		t := time.Unix(data.Timestamp, 0)
		data.ParsedTime = &t
	}

	for _, data := range w.DailyWeather {
		t := time.Unix(data.Timestamp, 0)
		data.ParsedTime = &t
	}

	return nil
}
