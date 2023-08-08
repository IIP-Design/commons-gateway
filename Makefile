.PHONY: build clean deploy

DEV_DIR = .gateway-dev
DEV_ENV  = -e DB_HOST=host.docker.internal:5454 -e DB_NAME=gateway_dev?sslmode=disable -e DB_PASSWORD=gateway_dev -e DB_USER=gateway_dev -e JWT_SECRET=2fweb3m$ndj

# Simulated events
EVENT_ADMIN_NEW = ./events/admin-new.json
EVENT_GET_CREDS = ./events/get-creds.json
EVENT_GUEST_AUTH = ./events/guest-auth.json
EVENT_PROVISION = ./events/provision.json

build:
	cd serverless;\
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/admin-new admin-new/*.go;\
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/guest-auth guest-auth/*.go;\
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/get-creds get-creds/*.go;\
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/provision provision/*.go;

clean:
	cd serverless; rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	cd serverless; sls deploy --verbose

dev:
	cd $(DEV_DIR); docker-compose up -d

dev-reset:
	$(DEV_DIR)/scripts/reset.sh;

dev-reseed:
	$(DEV_DIR)/scripts/reset.sh;$(DEV_DIR)/scripts/seed.sh

local-provision: build
	cd serverless;\sls invoke local -f provision $(DEV_ENV) -p $(EVENT_PROVISION);

local-auth: build
	cd serverless;\sls invoke local -f guest-auth $(DEV_ENV) -p $(EVENT_GUEST_AUTH);

local-admin: build
	cd serverless;\sls invoke local -f admin-new $(DEV_ENV) -p $(EVENT_ADMIN_NEW);
	
local-creds: build
	cd serverless;\sls invoke local -f get-creds $(DEV_ENV) -p $(EVENT_GET_CREDS);