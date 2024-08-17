package main

import (
	"backend/config"
	"backend/server/db"
	"log"

	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var Env = config.Env.Main().InitConfig()

func main() {
	dbS, err := db.NewDatabase(&Env)
	if err != nil {
		log.Fatal(err)
	}

	driver, err := mysql.WithInstance(dbS.DB(), &mysql.Config{})
	if err != nil {
		log.Fatal(err)
	}
	m, err := db.Migrate("file://server/db/migrate", "mysql", driver)
	if err != nil {
		log.Fatal(err)
	}
	if err:=m.Force(20240706122843);err!=nil{ // dont forget to change force version
		log.Fatal(err)
	}
	log.Println("done force version")
}
