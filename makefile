APP_NAME = api
OUTPUT_DIR = bin

# Common build options
BUILD_OPTIONS = -ldflags="-s -w"

# Define the default target
.DEFAULT_GOAL := build

# Targets
ci:
	./scripts/ci.sh

bootstrap:
	./scripts/tools.sh
	go mod download

wire:
	wire ./...

test:
	go test -v ./...

build: wire
	go build -o $(OUTPUT_DIR)/$(APP_NAME) $(BUILD_OPTIONS) -v ./cmd/$(APP_NAME)

run: build
	$(OUTPUT_DIR)/$(APP_NAME)

compress:
	upx -9 --best $(OUTPUT_DIR)/$(APP_NAME)

release-v3: build-linux-v3 compress

release-v2: build-linux-v2 compress

release: build-linux compress

build-linux: wire
	GOOS=linux GOARCH=amd64 go build -o $(OUTPUT_DIR)/$(APP_NAME) $(BUILD_OPTIONS) -v ./cmd/$(APP_NAME)

build-linux-v3:
	GOAMD64=v3 $(MAKE) build-linux

build-linux-v2:
	GOAMD64=v2 $(MAKE) build-linux

clean:
	rm -rf bin/*

run-docker:
	# setup APP_VERSION env variable to the latest git tag version
	export APP_VERSION=$(git describe --tags --abbrev=0)
	docker-compose up -d

swagger:
	swag init -g ./cmd/api/main.go