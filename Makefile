# Default target: print this help message
.PHONY: help
.DEFAULT_GOAL := help
help:
	@echo 'Usage:'
	@echo '  make <target>'
	@echo ''
	@echo 'Targets:'
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/  /'

## genkit: Run Genkit locally
.PHONY: genkit
genkit:
	cd genkit && genkit start

## dev: Run http server locally
.PHONY: dev
dev:
	go run main.go

## deploy: Deploy Cloud Run
.PHONY: deploy
deploy:
	gcloud run deploy slack-bot-application \
		--port 3000 \
		--region=us-central1 \
		--source=. \
		--env-vars-file=.env.yaml \
		--allow-unauthenticated

## tidy: Tidy modfiles, format and lint .go files
.PHONY: tidy
tidy:
	go mod tidy -v
	go fmt ./...
	golangci-lint run
