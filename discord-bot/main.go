package main

import (
	"discord-bot/database"
	"discord-bot/set"
	"flag"
	"github.com/joho/godotenv"
	"log"
)

var (
	Token string
)

func init() {
	flag.StringVar(&Token, "t", "토큰값", "Bot Token")
	flag.Parse()
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading.env file")
	}

	db := database.Setup()
	set.Setup(db)
}
