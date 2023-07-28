.PHONY: build clean deploy

DEV_DIR = .gateway-dev
DEV_DB  = -e DB_HOST=host.docker.internal:5454 -e DB_NAME=gateway_dev?sslmode=disable -e DB_PASSWORD=gateway_dev -e DB_USER=gateway_dev

build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/provision provision/*.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

dev:
	cd $(DEV_DIR); docker-compose up -d

local: build
	sls invoke local -f provision $(DEV_DB)