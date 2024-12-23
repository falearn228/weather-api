package main

import (
	_ "embed"
	"log"
	"net/http"
	"os"

	// "project/internal/api"
	// "project/pkg/redis"
	"github.com/joho/godotenv"
	routes "github.com/weather-api/internal/routes"
	redis "github.com/weather-api/service"
)

//go:embed config.env
var configData []byte



func main() {
	envMap, err := godotenv.UnmarshalBytes(configData)
	if err != nil {
		log.Fatalf("Failed to parse .env file... %v", err)
	}
	setenv(envMap)
	
	redisClient := redis.NewRedisClient()

	router := routes.SetupRoutes(redisClient)


	// Запуск сервера
	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setenv(envMap map[string]string) {
	apiKey := envMap["WEATHER_API_KEY"]
	redisAddr := envMap["REDIS_ADDR"]
	redisPassword := envMap["REDIS_PASSWORD"]

	os.Setenv("WEATHER_API_KEY", apiKey)
	os.Setenv("REDIS_ADDR", redisAddr)
	os.Setenv("REDIS_PASSWORD", redisPassword)
}