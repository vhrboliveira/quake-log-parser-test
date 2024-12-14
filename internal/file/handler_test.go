package file

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vhrboliveira/quake-log-parser-test/internal/logparser"
)

func TestReadFile(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantErr  bool
		expected []string
	}{
		{
			name: "Valid log file",
			content: `0:00 InitGame: \sv_floodProtect\1
20:34 ClientUserinfoChanged: 2 n\Player1\t\0
20:37 ShutdownGame:`,
			wantErr: false,
			expected: []string{
				"0:00 InitGame: \\sv_floodProtect\\1",
				"20:34 ClientUserinfoChanged: 2 n\\Player1\\t\\0",
				"20:37 ShutdownGame:",
			},
		},
		{
			name:     "Empty file",
			content:  "",
			wantErr:  false,
			expected: []string(nil),
		},
		{
			name:     "Non-existent file",
			content:  "",
			wantErr:  true,
			expected: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var tmpfile *os.File
			var tmpfileName string

			if !tc.wantErr {
				// Create temporary file
				var err error
				tmpfile, err = os.CreateTemp("", "test-*.log")
				assert.NoError(t, err)
				tmpfileName = tmpfile.Name()
				defer os.Remove(tmpfileName)

				// Write content to file
				_, err = tmpfile.Write([]byte(tc.content))
				assert.NoError(t, err)
				tmpfile.Close()
			} else {
				tmpfileName = "/nonexistent/file.log"
			}

			// Test ReadFile
			lines := make(chan string)
			errChan := make(chan error, 1)
			done := make(chan bool)

			go func() {
				ReadFile(tmpfileName, lines, errChan)
				close(done)
			}()

			// Collect results
			var results []string
			resultsDone := make(chan bool)

			go func() {
				for line := range lines {
					results = append(results, line)
				}
				resultsDone <- true
			}()

			// Wait for completion or error
			select {
			case err := <-errChan:
				if tc.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			case <-resultsDone:
				if tc.wantErr {
					t.Error("Expected error but got none")
				} else {
					assert.Equal(t, tc.expected, results)
				}
			case <-time.After(time.Second):
				t.Error("Test timed out")
			}
		})
	}
}

func TestWriteFile(t *testing.T) {
	tests := []struct {
		name    string
		report  logparser.GameReport
		setup   func() (string, func())
		wantErr bool
	}{
		{
			name: "Valid game report",
			report: logparser.GameReport{
				{
					"game_1": logparser.MatchReport{
						TotalKills: 1,
						Players:    []string{"Player1"},
						Kills:      map[string]int{"Player1": 1},
						KillsByMeans: map[string]int{
							logparser.MOD_ROCKET: 1,
						},
					},
				},
			},
			setup: func() (string, func()) {
				tmpdir, err := os.MkdirTemp("", "test-*")
				assert.NoError(t, err)
				return filepath.Join(tmpdir, "test.log"), func() {
					os.RemoveAll(tmpdir)
				}
			},
			wantErr: false,
		},
		{
			name:   "Empty report",
			report: logparser.GameReport{},
			setup: func() (string, func()) {
				tmpdir, err := os.MkdirTemp("", "test-*")
				assert.NoError(t, err)
				return filepath.Join(tmpdir, "test.log"), func() {
					os.RemoveAll(tmpdir)
				}
			},
			wantErr: false,
		},
		{
			name: "Fail to create file - invalid directory",
			report: logparser.GameReport{
				{
					"game_1": logparser.MatchReport{
						TotalKills: 1,
					},
				},
			},
			setup: func() (string, func()) {
				return "/nonexistent/directory/test.log", func() {}
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testFile, cleanup := tc.setup()
			defer cleanup()

			done := make(chan bool)
			gameReport := make(chan logparser.GameReport)
			errChan := make(chan error, 1)

			go func() {
				gameReport <- tc.report
				close(gameReport)
			}()

			go WriteFile(testFile, gameReport, done, errChan)

			select {
			case err := <-errChan:
				if tc.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			case <-done:
				if tc.wantErr {
					t.Error("Expected error but got none")
				} else {
					content, err := os.ReadFile(testFile + ".json")
					assert.NoError(t, err)

					var report logparser.GameReport
					err = json.Unmarshal(content, &report)
					assert.NoError(t, err)
					assert.Equal(t, tc.report, report)
				}
			case <-time.After(time.Second):
				t.Error("Test timed out")
			}
		})
	}
}
