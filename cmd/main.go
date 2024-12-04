package main

import (
	"log"
	"music_library/config"
)





func main () {
	config := config.LoadConfig()
	log.Printf("Configuration loaded: %+v\n", config)
	
}