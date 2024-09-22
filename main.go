package main

import (
	"context"
	"flag"
	"log"

	tgClient "go.mod/clients/telegram"
	"go.mod/consumer/eventconsumer"
	"go.mod/events/telegram"
	"go.mod/storage/sqlite"
)

const (
	tgBotHost         = "api.telegram.org"
	sqliteStoragePath = "data/sqlite/storage.db"
	batchSize         = 100
)

func main() {

	s, err := sqlite.New(sqliteStoragePath)

	if err != nil {
		log.Fatal("cannot create sqlite storage", err)
	}

	if err := s.Init(context.TODO()); err != nil {
		log.Fatal("cannot init sqlite storage", err)
	}

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		s,
	)

	log.Print("service started")

	consumer := eventconsumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func mustToken() string {
	token := flag.String(
		"tg-bot-token",
		"",
		"token for access to telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}
	return *token
}
