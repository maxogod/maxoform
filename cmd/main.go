package main

import (
	"log"

	"github.com/maxogod/maxoform/config"
	"github.com/maxogod/maxoform/internal/app"
	"github.com/maxogod/maxoform/internal/logger"
)

// TODO: maxoform apply, maxoform update
func main() {
	if err := logger.Init(true); err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	flagCfg, err := config.LoadFlagConf()
	if err != nil {
		log.Fatal(err)
	}

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
