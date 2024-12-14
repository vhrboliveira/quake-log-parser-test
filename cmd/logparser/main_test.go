package main

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vhrboliveira/quake-log-parser-test/internal/logparser"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name       string
		envVar     string
		errMessage string
	}{
		{
			name:       "Error - Missing LOG_FILE environment variable",
			envVar:     "",
			errMessage: "environment variable LOG_FILE is not set",
		},
		{
			name:       "Error - Non-existent log file",
			envVar:     "../../assets/nonexistent.log",
			errMessage: "error processing the log file: failed to open quake log file:",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.envVar != "" {
				os.Setenv("LOG_FILE", tc.envVar)
				defer os.Unsetenv("LOG_FILE")
				defer os.Remove(tc.envVar + ".json")
			} else {
				os.Unsetenv("LOG_FILE")
			}

			err := run()

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.errMessage)
		})
	}
}

func TestMainE2E(t *testing.T) {
	testLogFile := "../../assets/test.log"
	outputFile := testLogFile + ".json"

	killsByMeans := map[string]int{
		logparser.MOD_UNKNOWN:        0,
		logparser.MOD_SHOTGUN:        0,
		logparser.MOD_GAUNTLET:       0,
		logparser.MOD_MACHINEGUN:     0,
		logparser.MOD_GRENADE:        0,
		logparser.MOD_GRENADE_SPLASH: 0,
		logparser.MOD_ROCKET:         0,
		logparser.MOD_ROCKET_SPLASH:  0,
		logparser.MOD_PLASMA:         0,
		logparser.MOD_PLASMA_SPLASH:  0,
		logparser.MOD_RAILGUN:        0,
		logparser.MOD_LIGHTNING:      0,
		logparser.MOD_BFG:            0,
		logparser.MOD_BFG_SPLASH:     0,
		logparser.MOD_WATER:          0,
		logparser.MOD_SLIME:          0,
		logparser.MOD_LAVA:           0,
		logparser.MOD_CRUSH:          0,
		logparser.MOD_TELEFRAG:       0,
		logparser.MOD_FALLING:        0,
		logparser.MOD_SUICIDE:        0,
		logparser.MOD_TARGET_LASER:   0,
		logparser.MOD_TRIGGER_HURT:   0,
		logparser.MOD_NAIL:           0,
		logparser.MOD_CHAINGUN:       0,
		logparser.MOD_PROXIMITY_MINE: 0,
		logparser.MOD_KAMIKAZE:       0,
		logparser.MOD_JUICED:         0,
		logparser.MOD_GRAPPLE:        0,
	}

	if _, err := os.Stat(outputFile); err == nil {
		backupFile := outputFile + ".bak"
		os.Rename(outputFile, backupFile)
		defer os.Rename(backupFile, outputFile)
	}

	os.Setenv("LOG_FILE", testLogFile)
	defer os.Unsetenv("LOG_FILE")

	run()

	content, err := os.ReadFile(outputFile)
	assert.NoError(t, err)

	var report logparser.GameReport
	err = json.Unmarshal(content, &report)
	assert.NoError(t, err)

	assert.Len(t, report, 3, "Expected 3 games in report")

	// Game 1 verification
	game1 := report[0]["game_1"]
	assert.NotNil(t, game1)
	assert.Equal(t, 0, game1.TotalKills)
	assert.Len(t, game1.Players, 1)
	assert.Equal(t, game1.KillsByMeans, killsByMeans)
	assert.Contains(t, game1.Players, "Isgalamido (ID 2)")

	// Verify kills in game 1
	assert.Equal(t, 0, game1.Kills["Isgalamido (ID 2)"])

	// Game 2 verification
	game2 := report[1]["game_2"]
	killsByMeansGame2 := killsByMeans
	killsByMeansGame2[logparser.MOD_ROCKET] = 1
	killsByMeansGame2[logparser.MOD_ROCKET_SPLASH] = 8
	killsByMeansGame2[logparser.MOD_TRIGGER_HURT] = 1
	assert.NotNil(t, game2)
	assert.Equal(t, 10, game2.TotalKills)
	assert.Len(t, game2.Players, 4)
	assert.Equal(t, game2.KillsByMeans, killsByMeansGame2)
	assert.Contains(t, game2.Players, "Isgalamido (ID 2)")
	assert.Contains(t, game2.Players, "Dono da Bola (ID 3)")
	assert.Contains(t, game2.Players, "Mocinha (ID 4)")
	assert.Contains(t, game2.Players, "Chessus (ID 6)")

	// Verify kills in game 2
	assert.Equal(t, 3, game2.Kills["Isgalamido (ID 2)"])
	assert.Equal(t, 1, game2.Kills["Dono da Bola (ID 3)"])
	assert.Equal(t, 1, game2.Kills["Mocinha (ID 4)"])
	assert.Equal(t, -1, game2.Kills["Chessus (ID 6)"])

	// Game 3 verification
	game3 := report[2]["game_3"]
	killsByMeansGame3 := killsByMeans
	killsByMeansGame3[logparser.MOD_ROCKET] = 3
	killsByMeansGame3[logparser.MOD_ROCKET_SPLASH] = 6
	killsByMeansGame3[logparser.MOD_TRIGGER_HURT] = 2
	assert.NotNil(t, game3)
	assert.Equal(t, 11, game3.TotalKills)
	assert.Len(t, game3.Players, 4)
	assert.Equal(t, game3.KillsByMeans, killsByMeansGame3)
	assert.Contains(t, game3.Players, "Isgalamido (ID 2)")
	assert.Contains(t, game3.Players, "Dono da Bola (ID 3)")
	assert.Contains(t, game3.Players, "Mocinha (ID 4)")
	assert.Contains(t, game3.Players, "Oootsimo (ID 5)")

	// Verify kills in game 3
	assert.Equal(t, 2, game3.Kills["Isgalamido (ID 2)"])
	assert.Equal(t, 4, game3.Kills["Oootsimo (ID 5)"])
	assert.Equal(t, 1, game3.Kills["Dono da Bola (ID 3)"])
	assert.Equal(t, -2, game3.Kills["Mocinha (ID 4)"])
}
