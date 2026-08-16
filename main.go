package main

import (
	"log"

	"gator/internal/cli"
)

func main() {
	if err := cli.Cmds.Run(&cli.Sta, cli.ConsoleCommand); err != nil {
		log.Fatalf("Error running command: %v", err)
	}
}
