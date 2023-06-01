package main

import (
	"flag"
	"userservice/cmd/server"
)

var cfg = flag.String("config", "config/config.yaml", "config file location")

// main main
func main() {
	flag.Parse()
	server.Run(*cfg)
}
