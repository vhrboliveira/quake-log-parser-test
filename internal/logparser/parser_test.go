package logparser

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

var defaultKillsByMeans = map[string]int{
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
}

func TestParseLines(t *testing.T) {
	tests := []struct {
		name     string
		lines    []string
		expected GameReport
	}{
		{
			name: "Simple game with one kill",
			lines: []string{
				"0:00 InitGame: \\sv_floodProtect\\1",
				"0:01 ClientUserinfoChanged: 2 n\\Player1\\t\\0",
				"0:01 ClientUserinfoChanged: 3 n\\Player2\\t\\0",
				"0:02 Kill: 2 3 7: Player1 killed Player2 by MOD_ROCKET",
				"0:03 ShutdownGame:",
			},
			expected: GameReport{
				{
					"game_1": MatchReport{
						TotalKills: 1,
						Players:    []string{"Player1 (ID 2)", "Player2 (ID 3)"},
						Kills: map[string]int{
							"Player1 (ID 2)": 1,
							"Player2 (ID 3)": 0,
						},
						KillsByMeans: func() map[string]int {
							kills := copyKillsByMeans(defaultKillsByMeans)
							kills[MOD_ROCKET] = 1
							return kills
						}(),
					},
				},
			},
		},
		{
			name: "Game with world kill and suicide",
			lines: []string{
				"0:00 InitGame: \\sv_floodProtect\\1",
				"0:01 ClientUserinfoChanged: 2 n\\Player1\\t\\0",
				"0:02 Kill: 1022 2 22: <world> killed Player1 by MOD_TRIGGER_HURT",
				"0:03 Kill: 2 2 7: Player1 killed Player1 by MOD_ROCKET",
				"0:04 ShutdownGame:",
			},
			expected: GameReport{
				{
					"game_1": MatchReport{
						TotalKills: 2,
						Players:    []string{"Player1 (ID 2)"},
						Kills: map[string]int{
							"Player1 (ID 2)": -2,
						},
						KillsByMeans: func() map[string]int {
							kills := copyKillsByMeans(defaultKillsByMeans)
							kills[MOD_TRIGGER_HURT] = 1
							kills[MOD_ROCKET] = 1
							return kills
						}(),
					},
				},
			},
		},
		{
			name: "Multiple games",
			lines: []string{
				"0:00 InitGame: \\sv_floodProtect\\1",
				"0:01 ClientUserinfoChanged: 2 n\\Player1\\t\\0",
				"0:01 ClientUserinfoChanged: 4 n\\Player2\\t\\0",
				"0:02 Kill: 2 4 7: Player1 killed Player2 by MOD_ROCKET",
				"0:03 ShutdownGame:",
				"0:04 InitGame: \\sv_floodProtect\\1",
				"0:05 ClientUserinfoChanged: 2 n\\Player1\\t\\0",
				"0:05 ClientUserinfoChanged: 4 n\\Player3\\t\\0",
				"0:06 Kill: 2 4 8: Player1 killed Player3 by MOD_ROCKET_SPLASH",
				"0:07 ShutdownGame:",
			},
			expected: GameReport{
				{
					"game_1": MatchReport{
						TotalKills: 1,
						Players:    []string{"Player1 (ID 2)", "Player2 (ID 4)"},
						Kills: map[string]int{
							"Player1 (ID 2)": 1,
							"Player2 (ID 4)": 0,
						},
						KillsByMeans: func() map[string]int {
							kills := copyKillsByMeans(defaultKillsByMeans)
							kills[MOD_ROCKET] = 1
							return kills
						}(),
					},
				},
				{
					"game_2": MatchReport{
						TotalKills: 1,
						Players:    []string{"Player1 (ID 2)", "Player3 (ID 4)"},
						Kills: map[string]int{
							"Player1 (ID 2)": 1,
							"Player3 (ID 4)": 0,
						},
						KillsByMeans: func() map[string]int {
							kills := copyKillsByMeans(defaultKillsByMeans)
							kills[MOD_ROCKET_SPLASH] = 1
							return kills
						}(),
					},
				},
			},
		},
		{
			name: "Two inits without shutdown",
			lines: []string{
				"0:00 InitGame: \\sv_floodProtect\\1",
				"0:01 ClientUserinfoChanged: 2 n\\Player1\\t\\0",
				"0:01 ClientUserinfoChanged: 3 n\\Player2\\t\\0",
				"0:02 Kill: 2 3 7: Player1 killed Player2 by MOD_ROCKET",
				"0:03 InitGame: \\sv_floodProtect\\1",
				"0:04 ClientUserinfoChanged: 2 n\\Player1\\t\\0",
				"0:01 ClientUserinfoChanged: 4 n\\Player3\\t\\0",
				"0:05 Kill: 2 4 7: Player1 killed Player3 by MOD_ROCKET_SPLASH",
				"0:06 ShutdownGame:",
			},
			expected: GameReport{
				{
					"game_1": MatchReport{
						TotalKills: 1,
						Players:    []string{"Player1 (ID 2)", "Player2 (ID 3)"},
						Kills: map[string]int{
							"Player1 (ID 2)": 1,
							"Player2 (ID 3)": 0,
						},
						KillsByMeans: func() map[string]int {
							kills := copyKillsByMeans(defaultKillsByMeans)
							kills[MOD_ROCKET] = 1
							return kills
						}(),
					},
				},
				{
					"game_2": MatchReport{
						TotalKills: 1,
						Players:    []string{"Player1 (ID 2)", "Player3 (ID 4)"},
						Kills: map[string]int{
							"Player1 (ID 2)": 1,
							"Player3 (ID 4)": 0,
						},
						KillsByMeans: func() map[string]int {
							kills := copyKillsByMeans(defaultKillsByMeans)
							kills[MOD_ROCKET_SPLASH] = 1
							return kills
						}(),
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			lines := make(chan string)
			gameReport := make(chan GameReport)

			go func() {
				for _, line := range tc.lines {
					lines <- line
				}
				close(lines)
			}()

			go ParseLines(lines, gameReport)

			result := <-gameReport

			// Sort players in both expected and result for consistent comparison
			for i := range result {
				for gameName := range result[i] {
					sort.Strings(result[i][gameName].Players)
				}
			}
			for i := range tc.expected {
				for gameName := range tc.expected[i] {
					sort.Strings(tc.expected[i][gameName].Players)
				}
			}

			assert.Equal(t, tc.expected, result)
		})
	}
}

func copyKillsByMeans(original map[string]int) map[string]int {
	copy := make(map[string]int)
	for k, v := range original {
		copy[k] = v
	}
	return copy
}
