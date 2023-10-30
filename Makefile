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
EVENT_TEAM_CREATE = ./config/sim-events/team-create.json
EVENT_UPLOAD_METADATA = ./config/sim-events/upload-metadata.json

build:
	cd web; npm i && npm run build;
	cd serverless;\
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/authorizer funcs/authorizer/*.go;\
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/admin-create funcs/admin-create/*.go;\
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/admin-deactivate funcs/admin-deactivate/*.go;\
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/admin-get funcs/admin-get/*.go;\
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/admin-update funcs/admin-update/*.go;\
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/admins-get funcs/admins-get/*.go;\
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/aprimo-create-record funcs/aprimo-create-record/*.go;\
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/aprimo-upload-file funcs/aprimo-upload-file/*.go;\
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/email-2fa funcs/email-2fa/*.go;\
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/guest-approve funcs/guest-approve/*.go;\
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/guest-auth funcs/guest-auth/*.go;\
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/guest-deactivate funcs/guest-deactivate/*.go;\
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/guest-get funcs/guest-get/*.go;\
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/guest-reauth funcs/guest-reauth/*.go;\
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/guest-unlock funcs/guest-unlock/*.go;\
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/guest-update funcs/guest-update/*.go;\
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/guests-get funcs/guests-get/*.go;\
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/guests-pending funcs/guests-pending/*.go;\
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/init-db funcs/init-db/*.go;\
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/creds-2fa funcs/creds-2fa/*.go;\
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/creds-2fa-clear funcs/creds-2fa-clear/*.go;\
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/creds-salt funcs/creds-salt/*.go;\
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/creds-propose funcs/creds-propose/*.go;\
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/creds-provision funcs/creds-provision/*.go;\
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/password-change funcs/password-change/*.go;\
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/password-reset funcs/password-reset/*.go;\
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/seed-db funcs/seed-db/*.go;\
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/team-create funcs/team-create/*.go;\
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/team-update funcs/team-update/*.go;\
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/teams-get funcs/teams-get/*.go;\
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/upload-metadata funcs/upload-metadata/*.go;\
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/upload-presigned-url funcs/upload-presigned-url/*.go;\
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/uploader-get funcs/uploader-get/*.go;\

clean:
	cd serverless; rm -rf ./bin ./vendor Gopkg.lock;\
	cd ../web; rm -rf dist;

deploy: clean build
	cd serverless; npm run sls -- deploy --stage $(STAGE) --verbose;\
	aws s3 cp ../web/dist s3://commons-gateway-$(STAGE)-site/web/ --recursive;

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

local-guests-pending: build
	cd serverless;\
	npm run sls -- invoke local -f guestsPending $(DEV_DB_ENV);

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

local-upload-metadata: build
	cd serverless;\
	npm run sls -- invoke local -f uploadMetadata $(DEV_AWS_ENV) -p $(EVENT_UPLOAD_METADATA);

local-upload-presigned-url: build
	cd serverless;\
	npm run sls -- invoke local -f uploadPresignedUrl $(DEV_ENV);

local-uploader-get: build
	cd serverless;\
	npm run sls -- invoke local -f uploaderGet $(DEV_ENV);
