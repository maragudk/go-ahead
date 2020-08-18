.PHONY: build cover start test

build:
	go build -o ahead cmd/ahead/*.go

cover:
	go tool cover -html=cover.out

start: build
	./ahead start

test:
	go test -coverprofile=cover.out -mod=readonly ./...
