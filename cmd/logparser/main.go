package main

func main() {
	filePath := os.Getenv("LOG_FILE")
	if filePath == "" {
		log.Fatal("environment variable LOG_FILE is not set. Please set the LOG_FILE environment variable to the path of the quake log file. Ex: \"assets/quake.log\"")
	}
}
