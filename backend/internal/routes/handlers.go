package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
	weather "github.com/weather-api/internal/api"
)

func SetupRoutes(red *redis.Client) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/weather", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GET params were:", r.URL.Query())

		city := r.URL.Query().Get("city")
		if city == "" {
			http.Error(w, "City parameter is required", http.StatusBadRequest)
			return
		}

		weather, err := weather.GetWeather(w, red, city)
		if err != nil {
			http.Error(w, "Error while get weather from 3rd party API...", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(weather)
		if err != nil {
			http.Error(w, "Failed to encode weather data", http.StatusInternalServerError)
			return
		}
	})
	return mux
}