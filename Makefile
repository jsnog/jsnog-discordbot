BINARY_NAME=jsnog-bot
VERSION=0.0.1
GOARCH=$(shell go env GOARCH)
GOOS=$(shell go env GOOS)

all: build

build:
	cd src && \
	go get && \
	go build -o ../bin/$(BINARY_NAME)-$(VERSION)_$(GOOS)_$(GOARCH) -ldflags="-s -w" -trimpath *.go

docker-build:
	docker build -t jsnog-bot:v${VERSION} ./ 

clean:
	@rm -rf bin/
