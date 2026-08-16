package cli

import (
	"errors"
	"fmt"
	"log"
	"os"

	"gator/internal/config"
)

type State struct {
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
}

func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Arguments) < 1 {
		return errors.New("the login handler expects a single argument: The username")
	}

	if err := s.Config.SetUser(cmd.Arguments[0]); err != nil {
		return err
	}
	fmt.Printf("User: %s has been set.", cmd.Arguments[0])
	return nil
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
