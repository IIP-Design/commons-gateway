.PHONY: build clean deploy

DEV_DIR = .gateway-dev
DEV_ENV  = -e DB_HOST=host.docker.internal:5454 -e DB_NAME=gateway_dev?sslmode=disable -e DB_PASSWORD=gateway_dev -e DB_USER=gateway_dev

# Simulated events
EVENT_PROVISION = ./events/provision.json

build:
	cd serverless; env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/provision provision/*.go

clean:
	cd serverless; rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	cd serverless; sls deploy --verbose

dev:
	cd $(DEV_DIR); docker-compose up -d

dev-reset:
	$(DEV_DIR)/reset/reset.sh

local: build
	cd serverless; sls invoke local -f provision $(DEV_ENV) -p $(EVENT_PROVISION)