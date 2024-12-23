package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"
	configs "github.com/weather-api/configs"
)

type WeatherResponse struct {
	CurrentConditions struct {
		Temp        float32 `json:"temp"`        // Температура
		Description string  `json:"conditions"` // Описание погоды
	} `json:"currentConditions"`
}

func (wData WeatherResponse) MarshalBinary() ([]byte, error) {
	return json.Marshal(wData)
}

func (wData WeatherResponse) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &wData)
}

func GetWeather(w http.ResponseWriter, red *redis.Client, city string) (*WeatherResponse, error)  {
	var ctx = context.Background()
	var WeatherResponse WeatherResponse

	if city == "" {
		http.Error(w, "City parameter is required", http.StatusBadRequest)
	}

	val, err := red.Get(ctx, "weather"+city).Result()
	if err == redis.Nil {
		apiURL := configs.LoadConfig().GetWeatherAPIURL(city)

		resp, err := http.Get(apiURL)

		if err != nil {
			http.Error(w, "Failed to fetch weather data from API", http.StatusInternalServerError)
			return nil, fmt.Errorf("failed to fetch weather data: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			http.Error(w, "Failed to fetch valid weather data", http.StatusInternalServerError)
			return nil, fmt.Errorf("weather API returned status code %d", resp.StatusCode)
		}

		err = json.NewDecoder(resp.Body).Decode(&WeatherResponse)
		if err != nil {
			http.Error(w, "Failed to decode weather API response", http.StatusInternalServerError)
			return nil, fmt.Errorf("failed to decode weather API response: %v", err)
		}

		WeatherResponseBytes, err := WeatherResponse.MarshalBinary()
		if err != nil {
			http.Error(w, "Failed to serialize weather data", http.StatusInternalServerError)
			return nil, fmt.Errorf("failed to serialize weather data: %v", err)
		}

		err = red.Set(ctx, "weather:"+city, WeatherResponseBytes, 0).Err()
		if err != nil {
			log.Printf("Redis Set Error: %v\n", err)
			http.Error(w, "Failed to cache weather data in Redis", http.StatusInternalServerError)
			return nil, fmt.Errorf("failed to cache weather data: %v", err)
		}
		
		} else if err != nil {
			// Handle Redis errors (other than key not found)
			log.Printf("Redis Get Error: %v\n", err)
			http.Error(w, "Failed to fetch weather data from Redis", http.StatusInternalServerError)
			return nil, fmt.Errorf("failed to fetch weather data from Redis: %v", err)
		} else {
			// If data is found in Redis, unmarshal it
			err = WeatherResponse.UnmarshalBinary([]byte(val))
			if err != nil {
				http.Error(w, "Failed to decode cached weather data", http.StatusInternalServerError)
				return nil, fmt.Errorf("failed to decode cached weather data: %v", err)
			}
		}
		return &WeatherResponse, nil
}