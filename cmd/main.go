package main

import (
	"flag"

	"divoc.primea.se/app/client"
	"divoc.primea.se/app/server"
	"divoc.primea.se/util"
)

func main() {
	mode := flag.String("mode", "client", "'client' or 'server'")
	flag.StringVar(&util.ServerAddress, "server", "", "address to server")
	flag.Parse()

	if *mode == "server" {
		server.StartServer()
	} else {
		client.StartClient()
	}
}
