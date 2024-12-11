file ?= assets/qgames.log

docker-run:
	@LOG_FILE=$(file) docker compose -f compose-dev.yaml up --build

docker-down:
	@LOG_FILE=$(file) docker compose -f compose-dev.yaml down  

.PHONY: docker-run docker-down