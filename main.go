package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"log"
	"mongo/db"

	"mongo/repository"

	"os"
)

func main() {
	app := fiber.New()

	//	Load env
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading the env file")
	}

	config := db.Config{
		User:     os.Getenv("MONOG_USER"),
		Password: os.Getenv("MONGO_PASSWORD"),
		Host:     os.Getenv("MONGO_HOST"),
		Port:     os.Getenv("MONGO_PORT"),
		DBName:   os.Getenv("MONGO_DBNAME"),
	}

	instance, err := db.Connect(&config)

	if err != nil {
		log.Fatal("Error Connecting to the database")
	}

	r := repository.Repository{
		Instance: instance,
	}
	r.SetupRoutes(app)

	app.Listen(":8000")
}
