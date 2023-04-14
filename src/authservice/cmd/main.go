package main

import (
	"flag"

	"authservice/cmd/server"
)

var cfg = flag.String("config", "config/config.yaml", "config file location")

// main main
func main() {
	flag.Parse()
	server.Run(*cfg)

	/*
		server, err := server.InitServer(*cfg)
		if err != nil {
			panic("failed.")
		}
		server.Run()
	*/
}
