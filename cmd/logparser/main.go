package main

import (
	"log"
	"os"

	"github.com/vhrboliveira/quake-log-parser-test/internal/file"
	"github.com/vhrboliveira/quake-log-parser-test/internal/logparser"
)

func main() {
	filePath := os.Getenv("LOG_FILE")
	if filePath == "" {
		log.Fatal("environment variable LOG_FILE is not set. Please set the LOG_FILE environment variable to the path of the quake log file. Ex: \"assets/quake.log\"")
	}

	lines := make(chan string)
	gameReport := make(chan logparser.GameReport)
	done := make(chan bool)

	go file.ReadFile(filePath, lines)
	go logparser.ParseLines(lines, gameReport)
	go file.WriteFile(filePath, gameReport, done)

	<-done
}
