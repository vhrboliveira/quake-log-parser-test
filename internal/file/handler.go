package file

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/vhrboliveira/quake-log-parser-test/internal/logparser"
)

func ReadFile(path string, lines chan<- string) {
	fmt.Println("opening quake log file:", path)

	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to open quake log file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines <- scanner.Text()
	}

	close(lines)

	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading file: %v", err)
	}
}

func WriteFile(path string, gameReport <-chan logparser.GameReport, done chan<- bool) {
	fileName := path + ".json"
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	for event := range gameReport {
		jsonData, err := json.MarshalIndent(event, "", "  ")
		if err != nil {
			fmt.Printf("Error marshaling JSON: %v\n", err)
			return
		}

		_, err = file.Write(jsonData)
		if err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
			return
		}
	}

	done <- true
}
