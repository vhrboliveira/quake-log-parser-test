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
		if eventType != "" && matches != nil {
			processEvent(eventType, matches, game)
		}
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

func processEvent(eventType string, matches []string, game *gameState) {
	switch eventType {
	case INIT_GAME:
		if game.gameStarted {
			game.endGame()
		}

		game.gameStarted = true
		game.initGame()
	case KILL:
		/*
			^.*Kill: (\d+) (\d+) (\d+): (.+) killed (.+) by (.+)$
			(\d+) = killerID
			(\d+) = victimID
			(\d+) not used
			(.+) = killerName
			(.+) = victimName
			(.+) = method
		*/
		killerID, _ := strconv.Atoi(matches[1])
		victimID, _ := strconv.Atoi(matches[2])
		killerName := matches[4]
		victimName := matches[5]
		method := matches[6]

		game.handlePlayerKill(killerName, victimName, killerID, victimID)
		game.handleKillsByMeans(method)

	case USER_INFO:
		/*
			(\d+) n\\([^\\]+)\\t
			(\d+) = playerID
			n\\([^\\]+)\\t = playerName (extracted from n\playerName\t)
		*/
		playerID, _ := strconv.Atoi(matches[1])
		playerName := matches[2]
		game.updateUserInfo(playerID, playerName)

	case END_GAME:
		game.gameStarted = false
		game.endGame()
	}
}

func (game *gameState) initGame() {
	game.players = make(map[int]*playerInfo)
	game.totalGames++
	game.matchReport = MatchReport{
		TotalKills: 0,
		Players:    make([]string, 0),
		Kills:      make(map[string]int),
		KillsByMeans: map[string]int{
			MOD_UNKNOWN:        0,
			MOD_SHOTGUN:        0,
			MOD_GAUNTLET:       0,
			MOD_MACHINEGUN:     0,
			MOD_GRENADE:        0,
			MOD_GRENADE_SPLASH: 0,
			MOD_ROCKET:         0,
			MOD_ROCKET_SPLASH:  0,
			MOD_PLASMA:         0,
			MOD_PLASMA_SPLASH:  0,
			MOD_RAILGUN:        0,
			MOD_LIGHTNING:      0,
			MOD_BFG:            0,
			MOD_BFG_SPLASH:     0,
			MOD_WATER:          0,
			MOD_SLIME:          0,
			MOD_LAVA:           0,
			MOD_CRUSH:          0,
			MOD_TELEFRAG:       0,
			MOD_FALLING:        0,
			MOD_SUICIDE:        0,
			MOD_TARGET_LASER:   0,
			MOD_TRIGGER_HURT:   0,
			MOD_NAIL:           0,
			MOD_CHAINGUN:       0,
			MOD_PROXIMITY_MINE: 0,
			MOD_KAMIKAZE:       0,
			MOD_JUICED:         0,
			MOD_GRAPPLE:        0,
		},
	}
}

func (game *gameState) endGame() {
	for ID, player := range game.players {
		playerName := fmt.Sprintf("%s (ID %d)", player.name, ID)
		if _, ok := game.matchReport.Kills[playerName]; !ok {
			game.matchReport.Kills[playerName] = 0
			game.matchReport.Players = append(game.matchReport.Players, playerName)
		}
		game.matchReport.Kills[playerName] += player.kills
	}

	gameName := fmt.Sprintf("game_%d", game.totalGames)
	report := make(map[string]MatchReport)
	report[gameName] = game.matchReport
	game.gameReport = append(game.gameReport, report)

	game.players = nil
}

func (game *gameState) handlePlayerKill(killerName, victimName string, killerID, victimID int) {
	if killerName != WORLD {
		if _, ok := game.players[killerID]; !ok {
			game.players[killerID] = &playerInfo{name: killerName, kills: 0}
		}
	}

	if _, ok := game.players[victimID]; !ok && victimName != WORLD {
		game.players[victimID] = &playerInfo{
			name:  victimName,
			kills: 0,
		}
	}

	if killerName == WORLD || killerID == victimID {
		game.players[victimID].kills -= 1
	} else if killerID != victimID {
		game.players[killerID].kills += 1
	}

	game.matchReport.TotalKills += 1
}

func (game *gameState) updateUserInfo(playerID int, playerName string) {
	if _, ok := game.players[playerID]; !ok {
		game.players[playerID] = &playerInfo{
			name:  playerName,
			kills: 0,
		}
	} else {
		game.players[playerID].name = playerName
	}
}

func (game *gameState) handleKillsByMeans(method string) {
	if _, ok := game.matchReport.KillsByMeans[method]; !ok {
		game.matchReport.KillsByMeans[method] = 0
	}
	game.matchReport.KillsByMeans[method] += 1
}
