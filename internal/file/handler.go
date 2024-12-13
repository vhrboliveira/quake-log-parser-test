package file

import (
	"bufio"
	"fmt"
	"log"
	"os"

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
