package main

import (
	"flag"

	"productservice/cmd/server"
)

var cfg = flag.String("config", "config/config.yaml", "config file location")

// main main
func main() {
	flag.Parse()
	server.Run(*cfg)
}
