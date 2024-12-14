# Quake Log Parser

A Go application that parses Quake 3 Arena game logs and generates detailed match statistics in JSON format.

## Overview

This project parses Quake 3 Arena server logs to extract and analyze game data, including:
- Kill statistics
- Player rankings
- Death causes
- Match summaries

## Features

- Parses multiple game matches from a single log file
- Tracks player kills and deaths
- Handles special cases like world kills and suicides
- Groups kills by death causes
- Supports player name changes during matches
- Outputs detailed JSON reports

### Special Rules

1. When `<world>` kills a player or a player commits suicide:
   - The player loses 1 kill point
   - The kill is counted in total_kills
   - `<world>` is not listed in players or kills

2. Player name changes:
   - The system tracks players by ID (if a player disconnects, a new connection with the previous ID is considered the same player, even if his name changes).
   - Name changes are handled automatically
   - Kill counts persist across name changes
   - Each game might have more than one player with the same name, thus, the players' names are shown with their ID on the reports.

3. Shutdown entry missing on the log:
   - In case a shutdown entry is missing in the log between 2 matches, if the parser finds another **InitGame**, it closes the previous match and starts a new one.   

## Input Format

The parser expects a Quake 3 Arena log file to be provided as an environment variable (the path + file name and extension). If not provided, it will fallback to the default file on the **assets** folder.
```bash
file=assets/qgames.log
```

## Output Format

The parser generates a JSON file with the following structure:
```json
[
  {
    "game_1": {
      "total_kills": 45,
      "players": ["Player 1 (ID 2)", "Player 2 (ID 3)"],
      "kills": {
        "Player 1 (ID 2)": 5,
        "Player 2 (ID 3)": 3
      },
      "kills_by_means": {
        "MOD_ROCKET": 5,
        "MOD_RAILGUN": 2,
        // ... other death causes
      }
    }
  }
]
```

## Output Location

The parser generates a JSON file with the same name as the input file plus `.json` extension. For example:
- Input: `assets/qgames.log`
- Output: `assets/qgames.log.json`

## Requirements

- [Go 1.23 or higher](https://go.dev/doc/install)
- [Docker v27.3.1 or latest](https://docs.docker.com/engine/install/)
- [Make](https://www.gnu.org/software/make/)

## Running the Project

1. Clone the repository:
```bash
$ git clone https://github.com/vhrboliveira/quake-log-parser-test/
$ cd quake-log-parser-test
```

2. Available commands:

### Using Custom Log File
The environment variable `LOG_FILE` is required. You can set it manually or pass a custom log file path. All run commands accept a custom log file path:
```bash
$ make run file=assets/custom.log
$ make docker-prod-run file=assets/custom.log
$ make docker-dev-run file=assets/custom.log
```

### Build and Run
```bash
# Build the binary
$ make build

# Clean generated files
$ make clean

# Run directly with Go
$ make run

# Run using the binary
$ make run-bin
```

### Docker Commands
**Disclaimer**: Docker maps the `assets` folder to the container. In order to run the production mode, your log must be inside the `assets` folder. 
```bash
# Run with Docker in production mode
$ make docker-prod-run

# Stop production container
$ make docker-prod-down

# Run with Docker in development mode (hot reload)
$ make docker-dev-run

# Stop development container
$ make docker-dev-down
```

### Testing
```bash
# Run tests with coverage report
$ make tests

# Run tests with verbose output
$ make tests-verbose

# Show coverage in browser
$ make show-coverage

# Show coverage by function in terminal
$ make show-coverage-func
```

## Project Structure
```
.
├── cmd/
│ └── logparser/ # Main application entry point
├── internal/
│ ├── file/ # File handling operations
│ └── logparser/ # Core parsing logic
├── assets/ # Log files and output
├── compose-dev.yaml # Development Docker compose configuration
├── compose-prod.yaml # Production Docker compose configuration
├── local.Dockerfile # Development Dockerfile
├── Dockerfile # Production Dockerfile
├── Makefile # Build and run command
└── .air.toml # Air hot reloading
```

## Development

The project uses Air (with Docker) for hot reloading during development. Any code changes will automatically trigger a rebuild and restart of the application.