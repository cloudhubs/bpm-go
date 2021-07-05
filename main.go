package main

import "bpm-go/api"

func main() {
	server := &api.Server{
		Port: 8085,
	}
	server.InitializeAndRun()
}
