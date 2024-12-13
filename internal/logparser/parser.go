package logparser

import (
	"fmt"
	"strconv"
)

func ParseLines(lines <-chan string, gameReport chan<- GameReport) {
	game := &gameState{
		totalGames:  0,
		gameStarted: false,
		players:     make(map[int]*playerInfo),
		matchReport: MatchReport{},
		gameReport:  make(GameReport, 0),
	}

	for line := range lines {
		eventType, matches := parseLogLine(line)
	}

	gameReport <- game.gameReport

	close(gameReport)
}

func parseLogLine(line string) (eventType string, matches []string) {
	for eventType, regex := range RegexPatterns {
		if matches := regex.FindStringSubmatch(line); matches != nil {
			return eventType, matches
		}
	}

	return "", nil
}
