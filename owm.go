package owm

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"time"
)

const (
	baseURL = "http://api.openweathermap.org/data/2.5/onecall?"
)

// WeatherData holds the weather information returned by the OWM API
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

// CurrentData represents the current weather conditions
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

// WeatherDescription summarizes the current weather for human presentation
type WeatherDescription struct {
	ID          uint16 `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// HourlyWeatherSlice represents one element in the hourly forecast
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

// DailyWeatherSlice represents one element in the daily forecast
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

	for i := range w.HourlyWeather {
		t := time.Unix(w.HourlyWeather[i].Timestamp, 0)
		w.HourlyWeather[i].ParsedTime = &t
	}

	for i := range w.DailyWeather {
		t := time.Unix(w.DailyWeather[i].Timestamp, 0)
		w.DailyWeather[i].ParsedTime = &t
	}

	return nil
}

func (weather WeatherData) String() string {
	s, _ := json.MarshalIndent(weather, "", "\t")
	return string(s)
}

// WeatherAt returns the weather closest to referenceTime
func (weather *WeatherData) WeatherAt(referenceTime time.Time) *HourlyWeatherSlice {
	for _, slice := range weather.HourlyWeather {
		if slice.ParsedTime == nil {
			continue
		}
		difference := referenceTime.Sub(*slice.ParsedTime)
		if math.Abs(difference.Minutes()) < 30 {
			return &slice
		}
	}

	return nil
}

func (weather WeatherData) TemperatureAt(referenceTime time.Time) float64 {
	entry := weather.WeatherAt(referenceTime)
	if entry == nil {
		return -1
	}
	return entry.Temperature
}

func (weather WeatherData) WeatherTill(referenceTime time.Time) []HourlyWeatherSlice {
	afterLast := -1

	for i := range weather.HourlyWeather {
		timeStamp := weather.HourlyWeather[i].ParsedTime
		if timeStamp == nil {
			continue
		}
		difference := referenceTime.Sub(*timeStamp)
		if difference.Minutes() < -30 {
			afterLast = i
			break
		}
	}

	if afterLast == -1 {
		return nil
	}

	return weather.HourlyWeather[0:afterLast]
}

// CumulativePrecipitationTill returns the cumulative precipitation from the beginning of the data range till referenceTime
func (weather WeatherData) CumulativePrecipitationTill(referenceTime time.Time) float64 {
	forecast := weather.WeatherTill(referenceTime)
	if forecast == nil {
		return -1
	}

	val := 0.0

	for _, item := range forecast {
		val += item.Rain.OneHour
		val += item.Snow.OneHour
	}

	return val
}

// AverageTemperatureTill returns the average temperature from the beginning of the data range till referenceTime
func (weather WeatherData) AverageTemperatureTill(referenceTime time.Time) float64 {
	forecast := weather.WeatherTill(referenceTime)
	if forecast == nil {
		return -1
	}

	val := 0.0

	for _, item := range forecast {
		val += item.Temperature
	}

	return val / float64(len(forecast))
}

// AverageFeelsLikeTill returns the average feels like temperature from the beginning of the data range till referenceTime
func (weather WeatherData) AverageFeelsLikeTill(referenceTime time.Time) float64 {
	forecast := weather.WeatherTill(referenceTime)
	if forecast == nil {
		return -1
	}

	val := 0.0

	for _, item := range forecast {
		val += item.FeelsLike
	}

	return val / float64(len(forecast))
}
