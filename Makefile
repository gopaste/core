APP_NAME := pastebin
DOCKER_COMPOSE_FILE := docker-compose.yml

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

help:
	@echo "Usage:"
	@echo "  make build       - Compiles the application"
	@echo "  make run         - Starts the application using Docker Compose"
	@echo "  make stop        - Stops the application and removes the containers"
	@echo "  make clean       - Removes the containers and volumes"
	@echo "  make help        - Shows this help message"
