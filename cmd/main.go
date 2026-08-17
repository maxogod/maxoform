package main

import (
	"log"

	"github.com/maxogod/maxoform/config"
	"github.com/maxogod/maxoform/internal/app"
	"github.com/maxogod/maxoform/internal/logger"
)

// TODO: maxoform apply, maxoform update
func main() {
	flagCfg, err := config.LoadFlagConf()
	if err != nil {
		log.Fatal(err)
	}

	if err := logger.Init(true, flagCfg.QuietAll); err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	cfg, err := config.LoadConfig(flagCfg.ConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	manifest, err := config.LoadDconfManifest(flagCfg.DconfManifestPath)
	if err != nil {
		log.Fatal(err)
	}

	serviceManifest, err := config.LoadServiceManifest(flagCfg.ServicesManifestPath)
	if err != nil {
		log.Fatal(err)
	}

	appRunner := app.New(cfg, flagCfg, manifest, serviceManifest)

	if err := appRunner.Run(); err != nil {
		log.Fatal(err)
	}
}
