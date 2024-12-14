package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/vhrboliveira/quake-log-parser-test/internal/file"
	"github.com/vhrboliveira/quake-log-parser-test/internal/logparser"
)

func run() error {
	filePath := os.Getenv("LOG_FILE")
	if filePath == "" {
		return errors.New("environment variable LOG_FILE is not set. Please set the LOG_FILE environment variable to the path of the quake log file. Ex: \"assets/quake.log\"")
	}

	lines := make(chan string)
	gameReport := make(chan logparser.GameReport)
	done := make(chan bool)
	errChan := make(chan error, 1)

	go file.ReadFile(filePath, lines, errChan)
	go logparser.ParseLines(lines, gameReport)
	go file.WriteFile(filePath, gameReport, done, errChan)

	select {
	case <-done:
		fmt.Println("log parsing completed successfully")
		return nil
	case err := <-errChan:
		close(errChan)
		return fmt.Errorf("error processing the log file: %v", err)
	}
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
