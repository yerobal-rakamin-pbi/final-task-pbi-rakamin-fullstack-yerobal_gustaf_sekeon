.PHONY: swag-install
swag-install:
	@go install github.com/swaggo/swag/cmd/swag@latest

.PHONY: swag-init
swag-init:
	@/bin/rm -rf ./docs
	@swag init -g ./app/main.go -o ./docs --parseInternal

.PHONY: build
build:
	@go build -o app ./app

.PHONY: run
run: swag-init build
	@./app/app