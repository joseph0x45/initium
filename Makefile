VERSION := $(shell git describe --tags --abbrev=0)
APP := initium

release:
	GOOS=linux GOARCH=amd64 \
		go build -tags release -ldflags "-X main.version=$(VERSION)" -o $(APP)
