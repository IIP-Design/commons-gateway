.PHONY: build clean deploy

DEV_DIR = .gateway-dev
DEV_ENV  = -e DB_HOST=host.docker.internal:5454 -e DB_NAME=gateway_dev?sslmode=disable -e DB_PASSWORD=gateway_dev -e DB_USER=gateway_dev -e JWT_SECRET=2fweb3m$ndj

# Simulated events
EVENT_ADMIN_CREATE = ./serverless/config/sim-events/admin-create.json
EVENT_CREDS_SALT = ./serverless/config/sim-events/creds-salt.json
EVENT_CREDS_PROVISION = ./serverless/config/sim-events/creds-provision.json
EVENT_GUEST_AUTH = ./serverless/config/sim-events/guest-auth.json

build:
	cd serverless;\
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/admin-create funcs/admin-create/*.go;\
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/guest-auth funcs/guest-auth/*.go;\
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/creds-salt funcs/creds-salt/*.go;\
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/creds-provision funcs/creds-provision/*.go;

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
	cd serverless;\sls invoke local -f creds-provision $(DEV_ENV) -p $(EVENT_CREDS_PROVISION);

local-auth: build
	cd serverless;\sls invoke local -f guest-auth $(DEV_ENV) -p $(EVENT_GUEST_AUTH);

local-admin: build
	cd serverless;\sls invoke local -f admin-create $(DEV_ENV) -p $(EVENT_ADMIN_CREATE);
	
local-salt: build
	cd serverless;\sls invoke local -f creds-salt $(DEV_ENV) -p $(EVENT_CREDS_SALT);