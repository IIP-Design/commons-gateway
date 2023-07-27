.PHONY: build clean deploy

DEV_DIR = .gateway-dev

build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/creds creds/*.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

dev:
	cd $(DEV_DIR); docker-compose up -d

local: build
	sls invoke local -f creds