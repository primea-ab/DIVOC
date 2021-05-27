package main

import (
	"flag"
	"fmt"

	"divoc.primea.se/app/client"
	"divoc.primea.se/app/server"
)

func main() {
	mode := flag.String("mode", "client", "'client' or 'server'")
	flag.Parse()

	fmt.Println(*mode)

	if *mode == "server" {
		server.StartServer()
	} else {
		client.StartClient()
	}
}
