file ?= assets/qgames.log

build:
	@go build -o main ./cmd/logparser

clean:
	@rm -f main coverage.out
	@rm -rf tmp

run:
	@LOG_FILE=$(file) go run ./cmd/logparser/main.go

run-bin:
	@make build && LOG_FILE=$(file) ./main	 	

tests:
	@go test ./... -coverprofile=coverage.out

tests-verbose:
	@go test ./... -v -cover -coverprofile=coverage.out

show-coverage: tests
	@go tool cover -html=coverage.out

show-coverage-func: tests
	@go tool cover -func=coverage.out

docker-prod-run:
	@LOG_FILE=$(file) docker compose -f compose-prod.yaml up --build

docker-prod-down:
	@LOG_FILE=$(file) docker compose -f compose-prod.yaml down

docker-dev-run:
	@LOG_FILE=$(file) docker compose -f compose-dev.yaml up --build

docker-dev-down:
	@LOG_FILE=$(file) docker compose -f compose-dev.yaml down  

.PHONY: build clean run run-bin tests tests-verbose show-coverage show-coverage-func docker-prod-run docker-prod-down docker-dev-run docker-dev-down