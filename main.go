package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/timly278/simplebank/api"
	db "github.com/timly278/simplebank/db/sqlc"
	"github.com/timly278/simplebank/util"
)


func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	dbconn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store := db.NewStore(dbconn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

}