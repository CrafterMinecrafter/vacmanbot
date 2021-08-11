package main

import (
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/nyan2d/vacmanbot/app"
	"github.com/nyan2d/vacmanbot/config"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	var configPath string
	flag.StringVar(&configPath, "c", "config.json", "config file")
	flag.Parse()

	cfg, err := config.LoadConfigFromFile(configPath)
	if err != nil {
		log.Panic(err)
	}

	application := app.NewApp(
		cfg.TelegramToken,
		cfg.WebhookEndpoint,
		cfg.CertificateFile,
		cfg.PrivateKeyFile,
		cfg.DatabasePath,
	)
	application.Start()
}
