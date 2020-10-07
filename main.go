package main

import (
	"github.com/joho/godotenv"
	"os"
)

func main() {
	godotenv.Load()
	a := App{}
	a.Initialize(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"))

	a.Run(os.Getenv("PORT"))
}
