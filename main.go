package main

import "github.com/EdwinRincon/browersfc-api/server"

func main() {
	server := server.NewServer()
	server.Start()
}
