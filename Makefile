APP_NAME := pastebin
DOCKER_COMPOSE_FILE := docker-compose.yml

# Comandos
.PHONY: build run stop clean

build:
	@echo "Building $(APP_NAME)..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) build

run:
	@echo "Running $(APP_NAME)..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) up

stop:
	@echo "Stopping $(APP_NAME)..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) down

clean:
	@echo "Cleaning up $(APP_NAME)..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) down -v

# Ajuda
help:
	@echo "Uso:"
	@echo "  make build       - Compila a aplicação"
	@echo "  make run         - Inicia a aplicação usando Docker Compose"
	@echo "  make stop        - Para a aplicação e remove os containers"
	@echo "  make clean       - Remove os containers e volumes"
	@echo "  make help        - Mostra esta mensagem de ajuda"
