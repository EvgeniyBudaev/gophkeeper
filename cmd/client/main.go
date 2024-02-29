package main

import (
	"github.com/EvgeniyBudaev/gophkeeper/internal/client/cli"
	"github.com/EvgeniyBudaev/gophkeeper/internal/server/logger"
	"log"
)

func main() {
	l, err := logger.NewLogger()
	if err != nil {
		log.Fatal(err)
	}
	if err := cli.Execute(); err != nil {
		l.Errorf("error: %v", err)
	}
}
