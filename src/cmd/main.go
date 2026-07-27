package main

import (
	"log"

	"github.com/maxogod/maxoform/config"
	"github.com/maxogod/maxoform/internal/logger"
)

func main() {
	if err := logger.Init(true); err != nil {
		log.Fatal(err)
	}

	if _, err := config.Load(""); err != nil {
		log.Fatal(err)
	}

	defer logger.Sync()
}
