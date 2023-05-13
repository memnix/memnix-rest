APP_NAME = api

OUTPUT_DIR = bin

ci:
	./scripts/ci.sh

test:
	go test -v ./...

build:
	go build -o $(OUTPUT_DIR)/$(APP_NAME) -v .

run:
	go run ./...

release:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(OUTPUT_DIR)/$(APP_NAME) -v .
	upx -9 --best $(OUTPUT_DIR)/$(APP_NAME)

clean:
	rm -rf bin/*

