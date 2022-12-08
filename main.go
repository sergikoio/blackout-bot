package main

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"

	"blackout-bot/internal/bot"
	"blackout-bot/internal/db"
	"blackout-bot/internal/schedule"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	database, err := db.InitDB()
	if err != nil {
		log.Error().Err(err).Send()
		panic(err)
	}

	defer func(database *db.DB) {
		err := database.Close()
		if err != nil {
			log.Error().Err(err).Send()
			panic(err)
		}
	}(database)

	group, err := strconv.Atoi(os.Getenv("GROUP"))
	if err != nil {
		panic(err)
	}
	sch, err := schedule.NewSchedule(group, "schedule.json")
	if err != nil {
		panic(err)
	}

	channelID, err := strconv.ParseInt(os.Getenv("CHANNEL_ID"), 10, 64)
	if err != nil {
		panic(err)
	}
	updMessageID, err := strconv.Atoi(os.Getenv("UPDATE_MESSAGE_ID"))
	if err != nil {
		panic(err)
	}

	log.Info().Msg("start init bot")

	err = bot.InitBot(os.Getenv("BOT_TOKEN"), channelID, updMessageID, database, sch)
	if err != nil {
		log.Error().Err(err).Send()
		panic(err)
	}
}
