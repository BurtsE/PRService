.PHONY: up down test help

## Start services (detached)
up:
	docker compose up -d

## Stop services
down:
	docker compose down


## Run tests with coverage
test:
	go test ./... --cover

## Show this help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | \
		sed 's/:.*?## /: /' | \
		sort