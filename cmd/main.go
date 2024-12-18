package main

import (
	"dating-site-api/server"
	"log"
)

func main() {
	Server := server.APIServer{}
	if err := Server.Run("8001"); err != nil {
		log.Fatalf("There are some issues during server starts")
	}
}
