MAIN_FILE := ./build/server.js

.PHONY: all
all: test run

.PHONY: test
test: start-ancillary test-migrate
	@echo "Running all tests"
	node ace test
	docker compose down


.PHONY: run
run: start-ancillary dev-migrate
	@echo "Starting application in debug mode..."
	node ace serve --watch


.PHONY: build
build:
	@echo "Generating application build for production...."
	node ace build --production


.PHONY: preview
preview: build start-ancillary dev-migrate
	@echo "Previewing application in production mode"
	node $(MAIN_FILE)


.PHONY: test-migrate
test-migrate:
	@echo "Running migrations in test environment..."
	NODE_ENV=test node ace migration:fresh


.PHONY: dev-migrate
dev-migrate: 
	@echo "Running migrations in development environment..."
	node ace migration:fresh


.PHONY: start-ancillary
start-ancillary:
	@echo "Starting database containers"
	docker compose up -d
	

.PHONY: stop-ancillary
stop-ancillary:
	@echo "Shutting down database containers"
	docker compose down 
	