package main

import (
	"backend/server/api"
	"log"
)

func main() {
	log.Fatal(api.NewApi().Run())
}
