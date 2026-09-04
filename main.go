package main

import (
	"database/sql"
	"log"
	"os"

	"gator/internal/config"
	"gator/internal/database"

	_ "github.com/lib/pq"
)

type State struct {
	Db     *database.Queries
	Config *config.Config
}

type Command struct {
	Name      string
	Arguments []string
}

type Commands struct {
	Cmds map[string]func(*State, Command) error
}

var (
	Sta            State
	Cmds           Commands
	ConsoleCommand Command
)

func main() {
	Init()
	db, err := sql.Open("postgres", config.Cfg.Db_url)
	dbQueries := database.New(db)
	Sta.Db = dbQueries
	if err != nil {
		log.Fatalf("Error connecting to the db: %v", err)
	}
	if err := Cmds.Run(&Sta, ConsoleCommand); err != nil {
		log.Fatalf("Error running command: %v", err)
	}
}

func Init() {
	args := os.Args
	if len(args) < 2 {
		log.Fatalf("Error: Not enough arguments")
		return
	}
	command_part := args[1]
	arguments_part := args[2:]

	Sta.Config = &config.Cfg

	Cmds = Commands{Cmds: make(map[string]func(*State, Command) error)}

	ConsoleCommand = Command{Name: command_part, Arguments: arguments_part}

	Cmds.Register("login", HandlerLogin)
	Cmds.Register("register", HandlerRegister)
	Cmds.Register("reset", HandlerReset)
	Cmds.Register("users", HandlerUsers)
	Cmds.Register("addfeed", middlewareLoggedIn(HandlerAddFeed))
	Cmds.Register("feeds", HandlerGetFeeds)
	Cmds.Register("follow", middlewareLoggedIn(HandlerFollow))
	Cmds.Register("following", middlewareLoggedIn(HandlerFollowing))
	Cmds.Register("unfollow", middlewareLoggedIn(HandlerUnfollow))
	Cmds.Register("agg", HandlerAgg)
	Cmds.Register("browse", HandlerBrowse)
}
