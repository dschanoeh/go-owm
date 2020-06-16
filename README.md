# go-owm
[![Travis (master)](https://travis-ci.com/dschanoeh/go-owm.svg?branch=master)](https://travis-ci.com/dschanoeh/go-owm)


go-owm is a library for the [OpenWeatherMap "One Call API"](https://openweathermap.org/api/one-call-api), providing current weather as well as hourly and daily forecasts.
See the [OWM API documentation](https://openweathermap.org/api/one-call-api) for a description of the provided elements.

## Usage

``` go
import (
    "fmt"
    "github.com/dschanoeh/go-owm"
    )

func main() {
    w, err := owm.GetWeather(52.41467, 10.74063, "OWM_API_KEY")

    fmt.printf("The current temperature is: %f", w.Current.Temperature)
}
```
