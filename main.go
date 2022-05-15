package main

import (
	"github.com/maciej-kapusta/gomongo/config"
	"github.com/maciej-kapusta/gomongo/web"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg := &config.Config{
		Port:     "3000",
		MongoDb:  "docs",
		MongoUri: "mongodb://localhost:27017",
	}
	server, err := web.SetupAll(cfg)
	if err != nil {
		fatal(err)
	}
	err = server.Run(":" + cfg.Port)
	if err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	log.Fatal().Msg(err.Error())
}
