COMPOSE := docker compose -f docker-compose.yaml

.PHONY: up down restart logs ps test build clean

up:
	$(COMPOSE) up --build -d

down:
	$(COMPOSE) down

restart:
	$(COMPOSE) down
	$(COMPOSE) up --build -d

logs:
	$(COMPOSE) logs -f

ps:
	$(COMPOSE) ps

test:
	cd backend && go test ./...

build:
	cd backend && go build ./...

clean:
	$(COMPOSE) down -v --remove-orphans
