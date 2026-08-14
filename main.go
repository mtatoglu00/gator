package main

import (
	"log"

	"gator/internal/config"
)

func main() {
	path, err := config.GetConfigFilePath()
	if err != nil {
		log.Println(err)
	}

	cfg := config.Read(path)
	cfg.SetUser("mtatoglu00")
}
