package main

import (
	"fmt"
	"log/slog"
	"os"
	"progression/discord"
	"progression/league"
	"progression/repository"
	"strconv"
)

type config struct {
	dcBotToken string
	pgDatabase string
	pgHostname string
	pgPassword string
	pgUsername string
	pgPort     int
}

func main() {
	conf := parseEnv()
	dataStore := repository.NewPostgresDataStore(
		conf.pgHostname,
		conf.pgPort,
		conf.pgUsername,
		conf.pgPassword,
		conf.pgDatabase,
	)
	err := dataStore.Connect()
	if err != nil {
		slog.Error("failed to connect to datastore", "error", err)
		return
	}
	leagueManager := league.NewLeagueManager(dataStore)
	discordBot, err := discord.New(conf.dcBotToken, leagueManager)
	if err != nil {
		slog.Error("failed to create discord bot", "error", err)
		return
	}

	err = discordBot.Start()
	slog.Info("discord bot ended", "error", err)
}

func parseEnv() config {
	conf := config{
		dcBotToken: os.Getenv("DC_BOT_TOKEN"),
		pgHostname: os.Getenv("PG_HOSTNAME"),
		pgDatabase: os.Getenv("PG_DATABASE"),
		pgUsername: os.Getenv("PG_USERNAME"),
		pgPassword: os.Getenv("PG_PASSWORD"),
	}

	port, err := strconv.Atoi(os.Getenv("PG_PORT"))
	if err != nil {
		panic(fmt.Sprintf("PG_PORT environment variable not set to a valid value: %v", err))
	}
	conf.pgPort = port

	return conf
}
