package main

import (
	"fmt"
	"os"

	"github.com/hreshchyshynt/gator/internal/config"
)

func main() {
	config, err := config.Read()
	if err != nil {
		fmt.Printf("error read config: %v", err)
		os.Exit(1)
	}

	fmt.Printf("Config read: %v", config)
}
