package main

import (
	"log"
	"reimagined-chainsaw/server"
)

func main() {
	log.Println("Initializing Server...")
	s := server.NewServer()
	s.WithHandlers()
	s.RunCallbackFn()
	s.Start()
	log.Println("Server initialized with success!")
}
