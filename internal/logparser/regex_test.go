package logparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegexPatterns(t *testing.T) {
	tests := []struct {
		name      string
		line      string
		eventType string
		want      bool
		matches   []string
	}{
		{
			name:      "Valid InitGame",
			line:      "0:00 InitGame: \\sv_floodProtect\\1\\sv_maxPing\\0",
			eventType: INIT_GAME,
			want:      true,
			matches: []string{
				"0:00 InitGame: \\sv_floodProtect\\1\\sv_maxPing\\0",
				"\\sv_floodProtect\\1\\sv_maxPing\\0",
			},
		},
		{
			name:      "Valid Kill",
			line:      "20:54 Kill: 2 3 22: Isgalamido killed Dono da Bola by MOD_TRIGGER_HURT",
			eventType: KILL,
			want:      true,
			matches: []string{
				"20:54 Kill: 2 3 22: Isgalamido killed Dono da Bola by MOD_TRIGGER_HURT",
				"2",
				"3",
				"22",
				"Isgalamido",
				"Dono da Bola",
				"MOD_TRIGGER_HURT",
			},
		},
		{
			name:      "Valid UserInfo",
			line:      "20:34 ClientUserinfoChanged: 2 n\\Isgalamido\\t\\0\\model\\uriel/zael",
			eventType: USER_INFO,
			want:      true,
			matches: []string{
				"20:34 ClientUserinfoChanged: 2 n\\Isgalamido\\t",
				"2",
				"Isgalamido",
			},
		},
		{
			name:      "Valid ShutdownGame",
			line:      "20:37 ShutdownGame:",
			eventType: END_GAME,
			want:      true,
			matches: []string{
				"20:37 ShutdownGame:",
				"",
			},
		},
		{
			name:      "Invalid Kill format",
			line:      "Kill: invalid format",
			eventType: KILL,
			want:      false,
			matches:   nil,
		},
		{
			name:      "Invalid UserInfo format",
			line:      "ClientUserinfoChanged: invalid\\format",
			eventType: USER_INFO,
			want:      false,
			matches:   nil,
		},
		{
			name:      "Kill with special characters in names",
			line:      "20:54 Kill: 2 3 22: Player!@#$% killed Player&*() by MOD_TRIGGER_HURT",
			eventType: KILL,
			want:      true,
			matches: []string{
				"20:54 Kill: 2 3 22: Player!@#$% killed Player&*() by MOD_TRIGGER_HURT",
				"2",
				"3",
				"22",
				"Player!@#$%",
				"Player&*()",
				"MOD_TRIGGER_HURT",
			},
		},
		{
			name:      "UserInfo with special characters",
			line:      "20:34 ClientUserinfoChanged: 2 n\\Player!@#$%\\t\\0",
			eventType: USER_INFO,
			want:      true,
			matches: []string{
				"20:34 ClientUserinfoChanged: 2 n\\Player!@#$%\\t",
				"2",
				"Player!@#$%",
			},
		},
		{
			name:      "World kill",
			line:      "20:54 Kill: 1022 2 22: <world> killed Isgalamido by MOD_TRIGGER_HURT",
			eventType: KILL,
			want:      true,
			matches: []string{
				"20:54 Kill: 1022 2 22: <world> killed Isgalamido by MOD_TRIGGER_HURT",
				"1022",
				"2",
				"22",
				"<world>",
				"Isgalamido",
				"MOD_TRIGGER_HURT",
			},
		},
		{
			name:      "Self kill",
			line:      "20:54 Kill: 2 2 22: Isgalamido killed Isgalamido by MOD_ROCKET_SPLASH",
			eventType: KILL,
			want:      true,
			matches: []string{
				"20:54 Kill: 2 2 22: Isgalamido killed Isgalamido by MOD_ROCKET_SPLASH",
				"2",
				"2",
				"22",
				"Isgalamido",
				"Isgalamido",
				"MOD_ROCKET_SPLASH",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pattern := RegexPatterns[tc.eventType]
			matches := pattern.FindStringSubmatch(tc.line)
			got := matches != nil

			assert.Equal(t, tc.want, got, "Pattern match status mismatch")

			if tc.want {
				assert.NotEmpty(t, matches, "No Matches")
				assert.Equal(t, tc.matches, matches, "Regex matches mismatch")
				assert.Len(t, matches, len(tc.matches), "Number of matches mismatch")
			} else {
				assert.Nil(t, matches, "Expected nil matches for invalid pattern")
			}
		})
	}
}
