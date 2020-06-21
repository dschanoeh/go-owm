package owm

import (
	"encoding/json"
	"math"
	"os"
	"strings"
	"testing"
	"time"
)

const (
	TestData = `{"lat":52.41,"lon":10.74,"timezone":"Europe/Berlin","timezone_offset":7200,"current":{"dt":1592737994,"sunrise":1592708069,"sunset":1592768594,"temp":23.69,"feels_like":21.02,"pressure":1020,"humidity":40,"dew_point":9.31,"uvi":7.01,"clouds":60,"visibility":10000,"wind_speed":3.6,"wind_deg":290,"weather":[{"id":803,"main":"Clouds","description":"broken clouds","icon":"04d"}]},"hourly":[{"dt":1592737200,"temp":23.69,"feels_like":21.65,"pressure":1020,"humidity":40,"rain": {"1h": 1.12},"dew_point":9.31,"clouds":60,"wind_speed":2.7,"wind_deg":302,"weather":[{"id":803,"main":"Clouds","description":"broken clouds","icon":"04d"}]},{"dt":1592740800,"temp":22.93,"feels_like":21.26,"pressure":1020,"humidity":43,"dew_point":9.7,"clouds":63,"wind_speed":2.32,"wind_deg":299,"rain": {"1h": 1.12},"weather":[{"id":803,"main":"Clouds","description":"broken clouds","icon":"04d"}]},{"dt":1592744400,"temp":23.05,"feels_like":21.68,"pressure":1019,"humidity":44,"dew_point":10.15,"clouds":18,"wind_speed":2.07,"wind_deg":292,"weather":[{"id":801,"main":"Clouds","description":"few clouds","icon":"02d"}]},{"dt":1592748000,"temp":23.19,"feels_like":22.02,"pressure":1019,"humidity":46,"dew_point":10.94,"clouds":8,"wind_speed":2.1,"wind_deg":275,"weather":[{"id":800,"main":"Clear","description":"clear sky","icon":"01d"}]},{"dt":1592751600,"temp":22.75,"feels_like":21.34,"pressure":1019,"humidity":49,"dew_point":11.49,"clouds":2,"wind_speed":2.67,"wind_deg":258,"weather":[{"id":800,"main":"Clear","description":"clear sky","icon":"01d"}]},{"dt":1592755200,"temp":22.17,"feels_like":20.82,"pressure":1018,"humidity":52,"dew_point":12.01,"clouds":2,"wind_speed":2.75,"wind_deg":252,"weather":[{"id":800,"main":"Clear","description":"clear sky","icon":"01d"}]},{"dt":1592758800,"temp":21.52,"feels_like":20.72,"pressure":1019,"humidity":60,"dew_point":13.51,"clouds":1,"wind_speed":2.67,"wind_deg":250,"weather":[{"id":800,"main":"Clear","description":"clear sky","icon":"01d"}]},{"dt":1592762400,"temp":19.88,"feels_like":19.56,"pressure":1019,"humidity":68,"dew_point":13.92,"clouds":3,"wind_speed":2.17,"wind_deg":249,"weather":[{"id":800,"main":"Clear","description":"clear sky","icon":"01d"}]},{"dt":1592766000,"temp":17.4,"feels_like":16.56,"pressure":1019,"humidity":73,"dew_point":12.67,"clouds":0,"wind_speed":2.31,"wind_deg":259,"weather":[{"id":800,"main":"Clear","description":"clear sky","icon":"01d"}]},{"dt":1592769600,"temp":15.62,"feels_like":14.77,"pressure":1020,"humidity":74,"dew_point":11.03,"clouds":0,"wind_speed":1.68,"wind_deg":279,"weather":[{"id":800,"main":"Clear","description":"clear sky","icon":"01n"}]},{"dt":1592773200,"temp":14.73,"feels_like":14.4,"pressure":1020,"humidity":76,"dew_point":10.56,"clouds":0,"wind_speed":0.75,"wind_deg":305,"weather":[{"id":800,"main":"Clear","description":"clear sky","icon":"01n"}]},{"dt":1592776800,"temp":13.92,"feels_like":13.6,"pressure":1020,"humidity":79,"dew_point":10.38,"clouds":0,"wind_speed":0.66,"wind_deg":260,"weather":[{"id":800,"main":"Clear","description":"clear sky","icon":"01n"}]},{"dt":1592780400,"temp":13.92,"feels_like":13.26,"pressure":1020,"humidity":78,"dew_point":10.18,"clouds":18,"wind_speed":1.06,"wind_deg":242,"weather":[{"id":801,"main":"Clouds","description":"few clouds","icon":"02n"}]},{"dt":1592784000,"temp":13.74,"feels_like":12.91,"pressure":1020,"humidity":76,"dew_point":9.7,"clouds":16,"wind_speed":1.09,"wind_deg":224,"weather":[{"id":801,"main":"Clouds","description":"few clouds","icon":"02n"}]},{"dt":1592787600,"temp":14.9,"feels_like":13.62,"pressure":1020,"humidity":70,"dew_point":9.5,"clouds":81,"wind_speed":1.7,"wind_deg":235,"weather":[{"id":803,"main":"Clouds","description":"broken clouds","icon":"04n"}]},{"dt":1592791200,"temp":15.69,"feels_like":14.1,"pressure":1020,"humidity":66,"dew_point":9.52,"clouds":90,"wind_speed":2.09,"wind_deg":237,"weather":[{"id":804,"main":"Clouds","description":"overcast clouds","icon":"04n"}]},{"dt":1592794800,"temp":15.06,"feels_like":13.52,"pressure":1020,"humidity":73,"dew_point":10.37,"clouds":87,"wind_speed":2.36,"wind_deg":242,"weather":[{"id":804,"main":"Clouds","description":"overcast clouds","icon":"04d"}]},{"dt":1592798400,"temp":15.41,"feels_like":13.78,"pressure":1020,"humidity":76,"dew_point":11.37,"clouds":91,"wind_speed":2.87,"wind_deg":267,"weather":[{"id":804,"main":"Clouds","description":"overcast clouds","icon":"04d"}]},{"dt":1592802000,"temp":15.97,"feels_like":14.24,"pressure":1021,"humidity":82,"dew_point":12.97,"clouds":84,"wind_speed":3.76,"wind_deg":288,"weather":[{"id":803,"main":"Clouds","description":"broken clouds","icon":"04d"}]},{"dt":1592805600,"temp":16.52,"feels_like":15.17,"pressure":1021,"humidity":82,"dew_point":13.59,"clouds":85,"wind_speed":3.46,"wind_deg":297,"weather":[{"id":804,"main":"Clouds","description":"overcast clouds","icon":"04d"}]},{"dt":1592809200,"temp":17.85,"feels_like":16.39,"pressure":1022,"humidity":77,"dew_point":13.79,"clouds":50,"wind_speed":3.78,"wind_deg":288,"weather":[{"id":802,"main":"Clouds","description":"scattered clouds","icon":"03d"}]},{"dt":1592812800,"temp":19.58,"feels_like":17.67,"pressure":1022,"humidity":69,"dew_point":13.87,"clouds":27,"wind_speed":4.41,"wind_deg":286,"weather":[{"id":802,"main":"Clouds","description":"scattered clouds","icon":"03d"}]},{"dt":1592816400,"temp":20.52,"feels_like":18.19,"pressure":1023,"humidity":62,"dew_point":13.17,"clouds":18,"wind_speed":4.66,"wind_deg":287,"weather":[{"id":801,"main":"Clouds","description":"few clouds","icon":"02d"}]},{"dt":1592820000,"temp":21.33,"feels_like":18.74,"pressure":1023,"humidity":57,"dew_point":12.48,"clouds":13,"wind_speed":4.78,"wind_deg":291,"weather":[{"id":801,"main":"Clouds","description":"few clouds","icon":"02d"}]},{"dt":1592823600,"temp":21.96,"feels_like":19.23,"pressure":1023,"humidity":54,"dew_point":12.42,"clouds":11,"wind_speed":4.88,"wind_deg":290,"weather":[{"id":801,"main":"Clouds","description":"few clouds","icon":"02d"}]},{"dt":1592827200,"temp":22.47,"feels_like":19.63,"pressure":1023,"humidity":52,"dew_point":12.31,"clouds":10,"wind_speed":4.99,"wind_deg":290,"weather":[{"id":800,"main":"Clear","description":"clear sky","icon":"01d"}]},{"dt":1592830800,"temp":22.79,"feels_like":20.02,"pressure":1023,"humidity":51,"dew_point":12.26,"clouds":0,"wind_speed":4.9,"wind_deg":293,"weather":[{"id":800,"main":"Clear","description":"clear sky","icon":"01d"}]},{"dt":1592834400,"temp":22.86,"feels_like":20.13,"pressure":1023,"humidity":51,"dew_point":12.35,"clouds":0,"wind_speed":4.86,"wind_deg":295,"weather":[{"id":800,"main":"Clear","description":"clear sky","icon":"01d"}]},{"dt":1592838000,"temp":22.17,"feels_like":19.68,"pressure":1023,"humidity":54,"dew_point":12.6,"clouds":14,"wind_speed":4.63,"wind_deg":300,"weather":[{"id":801,"main":"Clouds","description":"few clouds","icon":"02d"}]},{"dt":1592841600,"temp":21.51,"feels_like":19.68,"pressure":1023,"humidity":59,"dew_point":13.26,"clouds":22,"wind_speed":4.02,"wind_deg":309,"weather":[{"id":801,"main":"Clouds","description":"few clouds","icon":"02d"}]},{"dt":1592845200,"temp":20.84,"feels_like":19.36,"pressure":1023,"humidity":62,"dew_point":13.44,"clouds":19,"wind_speed":3.58,"wind_deg":318,"weather":[{"id":801,"main":"Clouds","description":"few clouds","icon":"02d"}]},{"dt":1592848800,"temp":19.31,"feels_like":18.24,"pressure":1024,"humidity":68,"dew_point":13.42,"clouds":16,"wind_speed":2.98,"wind_deg":324,"weather":[{"id":801,"main":"Clouds","description":"few clouds","icon":"02d"}]},{"dt":1592852400,"temp":16.81,"feels_like":15.4,"pressure":1024,"humidity":74,"dew_point":12.31,"clouds":1,"wind_speed":2.96,"wind_deg":326,"weather":[{"id":800,"main":"Clear","description":"clear sky","icon":"01d"}]},{"dt":1592856000,"temp":14.68,"feels_like":13.03,"pressure":1025,"humidity":77,"dew_point":10.77,"clouds":0,"wind_speed":2.7,"wind_deg":335,"weather":[{"id":800,"main":"Clear","description":"clear sky","icon":"01n"}]},{"dt":1592859600,"temp":13.45,"feels_like":11.83,"pressure":1026,"humidity":79,"dew_point":10.05,"clouds":0,"wind_speed":2.33,"wind_deg":347,"weather":[{"id":800,"main":"Clear","description":"clear sky","icon":"01n"}]},{"dt":1592863200,"temp":12.45,"feels_like":10.91,"pressure":1026,"humidity":81,"dew_point":9.43,"clouds":0,"wind_speed":1.99,"wind_deg":351,"weather":[{"id":800,"main":"Clear","description":"clear sky","icon":"01n"}]},{"dt":1592866800,"temp":11.55,"feels_like":10.19,"pressure":1026,"humidity":84,"dew_point":9,"clouds":0,"wind_speed":1.61,"wind_deg":347,"weather":[{"id":800,"main":"Clear","description":"clear sky","icon":"01n"}]},{"dt":1592870400,"temp":10.85,"feels_like":9.73,"pressure":1026,"humidity":85,"dew_point":8.49,"clouds":0,"wind_speed":1.08,"wind_deg":333,"weather":[{"id":800,"main":"Clear","description":"clear sky","icon":"01n"}]},{"dt":1592874000,"temp":10.29,"feels_like":9.14,"pressure":1026,"humidity":86,"dew_point":8.16,"clouds":0,"wind_speed":0.99,"wind_deg":310,"weather":[{"id":800,"main":"Clear","description":"clear sky","icon":"01n"}]},{"dt":1592877600,"temp":9.96,"feels_like":8.79,"pressure":1026,"humidity":87,"dew_point":7.96,"clouds":0,"wind_speed":0.97,"wind_deg":285,"weather":[{"id":800,"main":"Clear","description":"clear sky","icon":"01n"}]},{"dt":1592881200,"temp":9.76,"feels_like":8.43,"pressure":1026,"humidity":87,"dew_point":7.87,"clouds":0,"wind_speed":1.13,"wind_deg":284,"weather":[{"id":800,"main":"Clear","description":"clear sky","icon":"01d"}]},{"dt":1592884800,"temp":10.57,"feels_like":9.37,"pressure":1026,"humidity":86,"dew_point":8.42,"clouds":0,"wind_speed":1.17,"wind_deg":279,"weather":[{"id":800,"main":"Clear","description":"clear sky","icon":"01d"}]},{"dt":1592888400,"temp":12.72,"feels_like":11.65,"pressure":1026,"humidity":80,"dew_point":9.4,"clouds":0,"wind_speed":1.35,"wind_deg":273,"weather":[{"id":800,"main":"Clear","description":"clear sky","icon":"01d"}]},{"dt":1592892000,"temp":14.77,"feels_like":13.57,"pressure":1027,"humidity":71,"dew_point":9.75,"clouds":0,"wind_speed":1.61,"wind_deg":301,"weather":[{"id":800,"main":"Clear","description":"clear sky","icon":"01d"}]},{"dt":1592895600,"temp":16.49,"feels_like":15.33,"pressure":1027,"humidity":66,"dew_point":10.2,"clouds":0,"wind_speed":1.76,"wind_deg":319,"weather":[{"id":800,"main":"Clear","description":"clear sky","icon":"01d"}]},{"dt":1592899200,"temp":17.93,"feels_like":16.66,"pressure":1027,"humidity":61,"dew_point":10.55,"clouds":0,"wind_speed":1.99,"wind_deg":320,"weather":[{"id":800,"main":"Clear","description":"clear sky","icon":"01d"}]},{"dt":1592902800,"temp":19.26,"feels_like":17.84,"pressure":1027,"humidity":56,"dew_point":10.47,"clouds":0,"wind_speed":2.2,"wind_deg":314,"weather":[{"id":800,"main":"Clear","description":"clear sky","icon":"01d"}]},{"dt":1592906400,"temp":20.55,"feels_like":18.94,"pressure":1026,"humidity":50,"dew_point":10.02,"clouds":0,"wind_speed":2.28,"wind_deg":314,"weather":[{"id":800,"main":"Clear","description":"clear sky","icon":"01d"}]}],"daily":[{"dt":1592737200,"sunrise":1592708069,"sunset":1592768594,"temp":{"day":22.93,"min":13.74,"max":23.69,"night":13.74,"eve":20.12,"morn":23.69},"feels_like":{"day":21.26,"night":12.91,"eve":19.8,"morn":21.26},"pressure":1020,"humidity":43,"dew_point":9.7,"wind_speed":2.32,"wind_deg":299,"weather":[{"id":803,"main":"Clouds","description":"broken clouds","icon":"04d"}],"clouds":63,"uvi":7.01},{"dt":1592823600,"sunrise":1592794484,"sunset":1592855003,"temp":{"day":22.47,"min":10.85,"max":22.47,"night":10.85,"eve":19.31,"morn":16.52},"feels_like":{"day":19.63,"night":9.73,"eve":18.24,"morn":15.17},"pressure":1023,"humidity":52,"dew_point":12.31,"wind_speed":4.99,"wind_deg":290,"weather":[{"id":800,"main":"Clear","description":"clear sky","icon":"01d"}],"clouds":10,"uvi":6.98},{"dt":1592910000,"sunrise":1592880902,"sunset":1592941410,"temp":{"day":22.3,"min":12.4,"max":22.61,"night":12.4,"eve":20.08,"morn":14.77},"feels_like":{"day":20.77,"night":11.42,"eve":19.95,"morn":13.57},"pressure":1026,"humidity":44,"dew_point":9.57,"wind_speed":2.04,"wind_deg":340,"weather":[{"id":800,"main":"Clear","description":"clear sky","icon":"01d"}],"clouds":0,"uvi":7.52},{"dt":1592996400,"sunrise":1592967322,"sunset":1593027813,"temp":{"day":22.89,"min":15.08,"max":23.45,"night":15.08,"eve":21.75,"morn":16.03},"feels_like":{"day":21.68,"night":14.32,"eve":22.48,"morn":15.61},"pressure":1025,"humidity":55,"dew_point":13.54,"wind_speed":3.23,"wind_deg":80,"weather":[{"id":800,"main":"Clear","description":"clear sky","icon":"01d"}],"clouds":0,"uvi":7.48},{"dt":1593082800,"sunrise":1593053746,"sunset":1593114213,"temp":{"day":24.57,"min":15.88,"max":25.38,"night":15.88,"eve":22.48,"morn":17.28},"feels_like":{"day":23.77,"night":15.69,"eve":22.48,"morn":16.73},"pressure":1022,"humidity":62,"dew_point":17.06,"wind_speed":4.42,"wind_deg":74,"weather":[{"id":500,"main":"Rain","description":"light rain","icon":"10d"}],"clouds":5,"rain":0.98,"uvi":7.28},{"dt":1593169200,"sunrise":1593140173,"sunset":1593200610,"temp":{"day":22.04,"min":16.52,"max":24.2,"night":16.52,"eve":22.7,"morn":18.54},"feels_like":{"day":22.35,"night":16.32,"eve":24.42,"morn":18.7},"pressure":1015,"humidity":73,"dew_point":17,"wind_speed":2.94,"wind_deg":96,"weather":[{"id":500,"main":"Rain","description":"light rain","icon":"10d"}],"clouds":75,"uvi":7.5},{"dt":1593255600,"sunrise":1593226603,"sunset":1593287003,"temp":{"day":25.36,"min":16.45,"max":25.36,"night":16.45,"eve":21.82,"morn":18.88},"feels_like":{"day":24.86,"night":16.74,"eve":23.68,"morn":19.95},"pressure":1010,"humidity":68,"dew_point":19.11,"wind_speed":5.35,"wind_deg":258,"weather":[{"id":501,"main":"Rain","description":"moderate rain","icon":"10d"}],"clouds":62,"rain":5.68,"uvi":7.24},{"dt":1593342000,"sunrise":1593313036,"sunset":1593373394,"temp":{"day":20.79,"min":14.01,"max":21.51,"night":14.01,"eve":20.25,"morn":17},"feels_like":{"day":19.27,"night":12.31,"eve":19.5,"morn":15.84},"pressure":1013,"humidity":64,"dew_point":13.83,"wind_speed":3.85,"wind_deg":247,"weather":[{"id":500,"main":"Rain","description":"light rain","icon":"10d"}],"clouds":92,"rain":1.64,"uvi":7.23}]}`
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

func ParseTestData() (*WeatherData, error) {
	d := WeatherData{}

	if err := json.NewDecoder(strings.NewReader(TestData)).Decode(&d); err != nil {
		return nil, err
	}

	if err := parseData(&d); err != nil {
		return nil, err
	}

	return &d, nil
}
func TestParsing(t *testing.T) {

	weather, err := ParseTestData()
	if err != nil {
		t.Error(err)
	}

	if weather.Current.Humidity != 40 {
		t.Error()
	}

	if weather.Current.Weather[0].ID != 803 {
		t.Error()
	}

	timeStamp := weather.Current.ParsedTime
	reference, _ := time.Parse(time.RFC3339, "2020-06-21T11:13:14+00:00")
	if !timeStamp.Equal(reference) {
		t.Error()
	}
}

func TestWeatherAt(t *testing.T) {
	weather, err := ParseTestData()
	if err != nil {
		t.Error(err)
	}

	referenceTime, err := time.Parse(time.RFC3339, "2020-06-21T12:29:59+00:00")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	slice := weather.WeatherAt(referenceTime)
	if slice == nil {
		t.Error("Slice not found")
		t.FailNow()
	}

	if slice.Timestamp != 1592740800 {
		t.Error("WeatherAt returned the wrong weather slice")
	}
}

func TestWeatherAtOutOfRange(t *testing.T) {
	weather, err := ParseTestData()
	if err != nil {
		t.Error(err)
	}

	referenceTime, err := time.Parse(time.RFC3339, "2020-06-21T08:00:00+00:00")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	slice := weather.WeatherAt(referenceTime)
	if slice != nil {
		t.Error("Out of range test returned valid data")
		t.FailNow()
	}
}

func TestTemperatureAt(t *testing.T) {
	weather, err := ParseTestData()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	referenceTime, err := time.Parse(time.RFC3339, "2020-06-21T12:29:59+00:00")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	temperature := weather.TemperatureAt(referenceTime)
	if math.Abs(temperature-22.93) > 0.1 {
		t.Errorf("Wrong temperature returned: %f", temperature)
		t.FailNow()
	}
}

func TestWeatherTill(t *testing.T) {
	weather, err := ParseTestData()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	referenceTime, err := time.Parse(time.RFC3339, "2020-06-21T13:31:00+00:00")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	forecastData := weather.WeatherTill(referenceTime)
	if len(forecastData) != 4 {
		t.Errorf("Forecast data has incorrect length: %d", len(forecastData))
		t.FailNow()
	}
}

func TestCumulativePrecipitationTill(t *testing.T) {
	weather, err := ParseTestData()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	referenceTime, err := time.Parse(time.RFC3339, "2020-06-21T13:31:00+00:00")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	val := weather.CumulativePrecipitationTill(referenceTime)

	if math.Abs(val-2.24) > 0.1 {
		t.Error("Cumulative precipitation not matching")
	}
}

func TestAverageTemperatureTill(t *testing.T) {
	weather, err := ParseTestData()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	referenceTime, err := time.Parse(time.RFC3339, "2020-06-21T13:31:00+00:00")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	val := weather.AverageTemperatureTill(referenceTime)

	if math.Abs(val-23.215) > 0.01 {
		t.Error("Average temperature not matching")
	}
}

func TestAverageFeelsLikeTill(t *testing.T) {
	weather, err := ParseTestData()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	referenceTime, err := time.Parse(time.RFC3339, "2020-06-21T13:31:00+00:00")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	val := weather.AverageFeelsLikeTill(referenceTime)

	if math.Abs(val-21.6525) > 0.01 {
		t.Error("Average feels like not matching")
	}
}
