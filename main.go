package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/jukkakansanaho/gator/internal/config"
)

func main() {
	currentUserName := "jukka"

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	if err := cfg.SetUser(currentUserName); err != nil {
		log.Fatalf("error setting user: %v", err)
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		log.Fatalf("error marshalling config: %v", err)
	}
	fmt.Println(string(data))
}

