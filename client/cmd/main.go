package main

import (
	"github.com/EvgeniyBudaev/gophkeeper/client/internal/cmd"
	"github.com/EvgeniyBudaev/gophkeeper/internal/logger"
	"log"
)

func main() {
	l, err := logger.NewLogger()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Execute(); err != nil {
		l.Errorf("error: %v", err)
	}
}
