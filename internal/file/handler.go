package file

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/vhrboliveira/quake-log-parser-test/internal/logparser"
)

func ReadFile(path string, lines chan<- string, errChan chan<- error) {
	fmt.Println("opening quake log file:", path)

	file, err := os.Open(path)
	if err != nil {
		errChan <- fmt.Errorf("failed to open quake log file: %w", err)
		close(lines)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines <- scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		errChan <- fmt.Errorf("error reading file: %w", err)
	}

	close(lines)
}

func WriteFile(path string, gameReport <-chan logparser.GameReport, done chan<- bool, errChan chan<- error) {
	fileName := path + ".json"
	file, err := os.Create(fileName)
	if err != nil {
		errChan <- fmt.Errorf("error creating file: %w", err)
		done <- false
		return
	}
	defer file.Close()

	for event := range gameReport {
		jsonData, err := json.MarshalIndent(event, "", "  ")
		if err != nil {
			errChan <- fmt.Errorf("error marshaling JSON: %w", err)
			done <- false
			return
		}

		_, err = file.Write(jsonData)
		if err != nil {
			errChan <- fmt.Errorf("error writing to file: %w", err)
			done <- false
			return
		}
	}

	done <- true
}
