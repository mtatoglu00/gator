package main

import (
	"database/sql"
	"log"

	"gator/internal/cli"
	"gator/internal/config"
	"gator/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", config.Cfg.Db_url)
	dbQueries := database.New(db)
	cli.Sta.Db = dbQueries
	if err != nil {
		log.Fatalf("Error connecting to the db: %v", err)
	}
	if err := cli.Cmds.Run(&cli.Sta, cli.ConsoleCommand); err != nil {
		log.Fatalf("Error running command: %v", err)
	}
}
