# Go parameters
GO := /usr/local/go/bin/go
GOFMT := /usr/local/go/bin/gofmt
GOTEST := /usr/local/go/bin/go test
BINARY_NAME := yenexpress-backend
BINARY_FILE := ./yenexpress-backend
TEST_DIR := ./test
TEST_SUBDIRS := $(wildcard $(TEST_DIR)/*)
TEST_TARGETS := $(addsuffix -test ,$(TEST_SUBDIRS))
MAIN_FILE := ./cmd/main.go

.PHONY: all
all: clean fmt test preview

.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)

.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@$(GOFMT) -w .

.PHONY: test
test:
	go test ./test/... -v

.PHONY: build
build: clean
	@echo "Building binary..."
	@$(GO) build -o $(BINARY_NAME) $(MAIN_FILE)

.PHONY: run
run: start-ancillary
	@echo "Running application..."
	@$(GO) run $(MAIN_FILE)


.PHONY: preview
preview: build start-ancillary
	@echo "Previewing application..."
	@$(BINARY_FILE)



.PHONY: start-ancillary
start-ancillary:
	@echo "Starting database containers"
	sudo docker compose up -d
	

.PHONY: stop-ancillary
stop-ancillary:
	@echo "Shutting down database containers"
	sudo docker compose down 
	