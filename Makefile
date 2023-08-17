.PHONY: build clean deploy

DEV_DIR = .gateway-dev
DEV_AWS_ENV = -e AWS_SES_REGION=us-east-1 -e SOURCE_EMAIL_ADDRESS=contentcommons@state.gov
DEV_DB_ENV  = -e DB_HOST=host.docker.internal:5454 -e DB_NAME=gateway_dev?sslmode=disable -e DB_PASSWORD=gateway_dev -e DB_USER=gateway_dev -e JWT_SECRET=2fweb3m$ndj

# Add SERVERLESS_STAGE=<stage> to make command to deploy to different stage
SERVERLESS_STAGE=dev

# Simulated events
EVENT_ADMIN_CREATE = ./config/sim-events/admin-create.json
EVENT_CREDS_SALT = ./config/sim-events/creds-salt.json
EVENT_CREDS_PROVISION = ./config/sim-events/creds-provision.json
EVENT_EMAIL_2FA = ./config/sim-events/email-2fa.json
EVENT_EMAIL_CREDS = ./config/sim-events/email-creds.json
EVENT_EMAIL_SUPPORT_STAFF = ./config/sim-events/email-support-staff.json
EVENT_GUEST_AUTH = ./config/sim-events/guest-auth.json

build:
	cd serverless;\
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/admin-create funcs/admin-create/*.go;\
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/guest-auth funcs/guest-auth/*.go;\
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/creds-salt funcs/creds-salt/*.go;\
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/creds-provision funcs/creds-provision/*.go;\
	cd funcs;\
	cd email-2fa && npm run zip && cd ../;\
	cd email-creds && npm run zip && cd ../;\
	cd email-support-staff && npm run zip && cd ../;

clean:
	cd serverless; rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	cd serverless; sls deploy --stage $(SERVERLESS_STAGE) --verbose

dev:
	cd $(DEV_DIR); docker-compose up -d

dev-reset:
	$(DEV_DIR)/scripts/reset.sh;

dev-reseed:
	$(DEV_DIR)/scripts/reset.sh;$(DEV_DIR)/scripts/seed.sh

local-provision: build
	cd serverless;\sls invoke local -f creds-provision $(DEV_DB_ENV) -p $(EVENT_CREDS_PROVISION);

local-auth: build
	cd serverless;\sls invoke local -f guest-auth $(DEV_DB_ENV) -p $(EVENT_GUEST_AUTH);

local-admin: build
	cd serverless;\sls invoke local -f admin-create $(DEV_DB_ENV) -p $(EVENT_ADMIN_CREATE);
	
local-salt: build
	cd serverless;\sls invoke local -f creds-salt $(DEV_DB_ENV) -p $(EVENT_CREDS_SALT);
	
local-email-2fa: build
	cd serverless;\sls invoke local -f email-2fa $(DEV_AWS_ENV) -p $(EVENT_EMAIL_2FA);
	
local-email-creds: build
	cd serverless;\sls invoke local -f email-creds $(DEV_AWS_ENV) -p $(EVENT_EMAIL_CREDS);

local-email-support-staff: build
	cd serverless;\sls invoke local -f email-support-staff $(DEV_ENV) -p $(EVENT_EMAIL_SUPPORT_STAFF);
