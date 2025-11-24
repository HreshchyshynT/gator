package main

import (
	"fmt"
	"log"

	cfg "github.com/hreshchyshynt/gator/internal/config"
)

func main() {
	config, err := cfg.Read()
	if err != nil {
		log.Fatalf("error read config: %v", err)
	}

	err = config.SetUser("lane")
	if err != nil {
		log.Fatalf("Failed to set user for config: %v", err)
	}
	config, err = cfg.Read()
	if err != nil {
		log.Fatalf("error read new config: %v", err)
	}
	fmt.Printf("New config read: %v", config)
}
