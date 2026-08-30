package cli

import (
	"log"
	"os"

	"gator/internal/config"
	"gator/internal/database"
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

func init() {
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
	Cmds.Register("agg", HandlerAgg)
	Cmds.Register("addfeed", HandlerAddFeed)
	Cmds.Register("feeds", HandlerGetFeeds)
}

func (c *Commands) Run(s *State, cmd Command) error {
	if err := c.Cmds[cmd.Name](s, cmd); err != nil {
		return err
	}
	return nil
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.Cmds[name] = f
}
