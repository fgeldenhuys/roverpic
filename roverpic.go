package main

import (
	"log"
	"roverpic/downloader"
	"roverpic/roverapi"
	"roverpic/server"
)

func main() {
	conf, err := DecodeConfig("roverpic.toml")
	if err != nil {
		log.Fatalf("Unable to read config file: %s", err)
	}
	if err := conf.Validate(); err != nil {
		log.Fatalf("Invalid config: %s", err)
	}

	dl := downloader.Init(conf.Downloader)
	roverAPI := roverapi.Init(conf.RoverAPI)
	srv := server.Init(conf.Server, roverAPI, dl)
	srv.ListenAndServe()
}
