package main

import (
	"backend/config"
	"backend/server/db"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

//  migrate create -ext sql -dir server/db/migrate add_users_table

var Env = config.Env.Main().InitConfig()

func main() {

	db, err := db.NewDatabase(&Env)
	if err != nil {
		log.Fatal(err)
	}

	driver, err := mysql.WithInstance(db.DB(), &mysql.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://server/db/migrate", "mysql", driver)
	if err != nil {
		log.Fatal(err)
	}

	cmd := os.Args[len(os.Args)-1]
	if cmd == "up" {
		err := m.Up()
		if err != nil {
			log.Fatal("gagal membuat :", err)
		}
		fmt.Println("proses up selesai")
	}
	if cmd == "down" {
		err := m.Down()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("proses down selesai")

	}

}
