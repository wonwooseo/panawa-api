PKG=github.com/wonwooseo/panawa-api
APP=api
VERSION=$(shell git describe --tags --dirty --always)

build:
	go build -o $(APP) -ldflags="-X $(PKG)/build.Version=$(VERSION) -X $(PKG)/build.BuildTime=$(shell date -Iseconds)" ./cmd

image:
	GOOS=linux GOARCH=amd64 go build -o $(APP)-deploy -ldflags="-X $(PKG)/build.Version=$(VERSION) -X $(PKG)/build.BuildTime=$(shell date -Iseconds) -w -s" ./cmd
	docker build -t wonwooseo/panawa-api:latest .
	docker scout cves wonwooseo/panawa-api:latest

clean:
	@-rm $(APP)

# run api server locally
run: api
	./api --config=sample-config.yaml --cors_allow_all

.PHONY: build
