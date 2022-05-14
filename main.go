package main

import (
	"github.com/maciej-kapusta/gomongo/config"
	"github.com/maciej-kapusta/gomongo/web"
	"github.com/rs/zerolog/log"
)

func main() {

	server, err := web.New(&config.Config{
		Port:     "3000",
		MongoDb:  "docs",
		MongoUri: "mongodb://localhost:27017",
	})
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	err = server.Serve()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}
