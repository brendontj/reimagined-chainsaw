package main

import (
	"reimagined-chainsaw/server"
)

func main() {
	s := server.NewServer()
	s.WithHandlers()
	s.SetHandleMessage()
	s.RunCallbackFn()
	s.Start()
}
