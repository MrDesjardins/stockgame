package util

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetDBEnv() (isProduction bool, dbHost, dbPort, dbUser, dbPassword, dbName, dbUrl string) {
	env := os.Getenv("GO_ENV")
	println("env", env)
	isProduction = env == "production"
	if !isProduction {
		err := godotenv.Load(".env")
		if err != nil {
			panic("Error loading .env file")
		}
	}
	dbHost = os.Getenv("DB_HOST")
	dbPort = os.Getenv("DB_PORT")
	dbUser = os.Getenv("DB_USER")
	dbPassword = os.Getenv("DB_PASSWORD")
	dbName = os.Getenv("DB_NAME")
	dbUrl = os.Getenv("DATABASE_URL")

	return
}

func GetWebEnv() (port string) {
	port = os.Getenv("VITE_API_PORT")
	if port == "" {
		fmt.Println("Please set the database environment variables: VITE_API_PORT")
		log.Fatal("Missing database environment variables")
	}
	return
}
