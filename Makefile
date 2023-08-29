.PHONY: build clean deploy

# Sets the envrionmental variables needed to test functions locally
DEV_DIR = .gateway-dev
DEV_AWS_ENV = -e AWS_SES_REGION=us-east-1 -e SOURCE_EMAIL_ADDRESS=contentcommons@state.gov
DEV_DB_ENV  = -e DB_HOST=host.docker.internal:5454 -e DB_NAME=gateway_dev -e DB_PASSWORD=gateway_dev -e DB_USER=gateway_dev -e JWT_SECRET=2fweb3m$ndj

# Sets the stage for the serverless deployment.
# Can be overridden in the CLI as so: `make target STAGE=mystage`
STAGE=dev

# Simulated events
EVENT_ADMIN_CREATE = ./config/sim-events/admin-create.json
EVENT_CREDS_SALT = ./config/sim-events/creds-salt.json
EVENT_CREDS_PROVISION = ./config/sim-events/creds-provision.json
EVENT_GUEST_AUTH = ./config/sim-events/guest-auth.json
EVENT_GUESTS_GET = ./config/sim-events/guests-get.json
EVENT_EMAIL_2FA = ./config/sim-events/email-2fa.json
EVENT_EMAIL_CREDS = ./config/sim-events/email-creds.json
EVENT_EMAIL_SUPPORT_STAFF = ./config/sim-events/email-support-staff.json
EVENT_TEAM_CREATE = ./config/sim-events/team-create.json

build:
	cd web; npm run build;
	cd serverless;\
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/admin-create funcs/admin-create/*.go;\
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/admins-get funcs/admins-get/*.go;\
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/admin-get funcs/admin-get/*.go;\
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/guest-auth funcs/guest-auth/*.go;\
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/guests-get funcs/guests-get/*.go;\
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/init-db funcs/init-db/*.go;\
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/creds-salt funcs/creds-salt/*.go;\
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/creds-provision funcs/creds-provision/*.go;\
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/team-create funcs/team-create/*.go;\
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/teams-get funcs/teams-get/*.go;\
	cd funcs;\
	cd email-2fa && npm run zip && cd ../;\
	cd email-creds && npm run zip && cd ../;\
	cd email-support-staff && npm run zip && cd ../;

clean:
	cd serverless; rm -rf ./bin ./vendor Gopkg.lock;\
	cd ../web; rm -rf dist;

deploy: clean build
	aws s3 cp ./web/dist s3://$(STAGE).gateway.gpalab.digital/ --recursive;\
	cd serverless; npm run sls -- deploy --stage $(STAGE) --verbose;

dev:
	cd $(DEV_DIR); docker-compose up -d

dev-reset:
	$(DEV_DIR)/scripts/reset.sh;

dev-reseed:
	$(DEV_DIR)/scripts/reset.sh;$(DEV_DIR)/scripts/seed.sh

init-db:
	cd serverless;\
	npm run sls -- invoke -f initDB --stage $(STAGE)

local-provision: build
	cd serverless;\
	npm run sls -- invoke local -f credsProvision $(DEV_DB_ENV) -p $(EVENT_CREDS_PROVISION);

local-auth: build
	cd serverless;\
	npm run sls -- invoke local -f guestAuth $(DEV_DB_ENV) -p $(EVENT_GUEST_AUTH);

local-admin: build
	cd serverless;\
	npm run sls -- invoke local -f adminCreate $(DEV_DB_ENV) -p $(EVENT_ADMIN_CREATE);

local-admins: build
	cd serverless;\
	npm run sls -- invoke local -f adminsGet $(DEV_DB_ENV);

local-guests: build
	cd serverless;\
	npm run sls -- invoke local -f guestsGet $(DEV_DB_ENV) -p $(EVENT_GUESTS_GET);

local-salt: build
	cd serverless;\
	npm run sls -- invoke local -f credsSalt $(DEV_DB_ENV) -p $(EVENT_CREDS_SALT);

local-team: build
	cd serverless;\
	npm run sls -- invoke local -f teamCreate $(DEV_DB_ENV) -p $(EVENT_TEAM_CREATE);

local-teams: build
	cd serverless;\
	npm run sls -- invoke local -f teamsGet $(DEV_DB_ENV);

local-email-2fa: build
	cd serverless;\
	npm run sls -- invoke local -f email2fa $(DEV_AWS_ENV) -p $(EVENT_EMAIL_2FA);

local-email-creds: build
	cd serverless;\
	npm run sls -- invoke local -f emailCreds $(DEV_AWS_ENV) -p $(EVENT_EMAIL_CREDS);

local-email-support-staff: build
	cd serverless;\
	npm run sls -- invoke local -f emailSupportStaff $(DEV_ENV) -p $(EVENT_EMAIL_SUPPORT_STAFF);
