BINARY := build/rogue_game
LOG_FILE := ./debug.log
JSON_DIR := ./datalayer

.PHONY: build run test update log clean

build:
	@go build -o $(BINARY) ./cmd/rogue_game

$(JSON_DIR):
	mkdir -p $@

run: build $(JSON_DIR)
	@$(BINARY)

test:
	@go test -v ./...

update:
	@git checkout develop
	@git pull --ff-only origin develop
	@echo "✅ develop updated"

log:
	cat $(LOG_FILE)
# add history clean
clean:
	@rm -rf $(BINARY)
	@rm -rf $(JSON_DIR)/*.json