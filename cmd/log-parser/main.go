package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	filePath := os.Getenv("LOG_FILE")
	if filePath == "" {
		log.Fatal("environment variable LOG_FILE is not set. Please set the LOG_FILE environment variable to the path of the quake log file. Ex: \"assets/quake.log\"")
	}

	fmt.Println("reading quake log file:", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed to open quake log file: %v", err)
	}
	defer file.Close()
}
