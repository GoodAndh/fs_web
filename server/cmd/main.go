package main

import (
	"backend/server/api"
	"log"
)

func main() {
	api := api.NewApi()
	err := api.Run()
	if err != nil {
		log.Fatal("fatal:", err)
	}
}
