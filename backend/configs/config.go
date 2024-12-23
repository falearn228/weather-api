package configs

import (
	"fmt"
	"os"
)

type Config struct {
	WeatherAPIKey string
	WeatherAPIURL string
	RedisAddr     string
	RedisPassword string
}

func LoadConfig() *Config {
	return &Config{
		WeatherAPIKey: os.Getenv("WEATHER_API_KEY"),
		WeatherAPIURL: "https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/%s?unitGroup=metric&key=%s&contentType=json",
		RedisAddr: os.Getenv("REDIS_ADDR"),
		// RedisPassword: os.Getenv("REDIS_PASSWORD"),
	}
}

func (c *Config) GetWeatherAPIURL(location string) string {
	return fmt.Sprintf(c.WeatherAPIURL, location, c.WeatherAPIKey)
}