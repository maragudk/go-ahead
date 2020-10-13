.PHONY: build build-css certs clean cover down lint	migrate-create migrate-down migrate-goto migrate-up start test test-down test-integration test-up up up-build

export NAME := ahead
export MIGRATE_DB_URL := "postgres://${NAME}:123@localhost:5432/${NAME}?sslmode=disable"

build:
	GOOS=linux GOARCH=amd64 go build -o ${NAME} cmd/server/*.go

build-css:
	NODE_ENV=production ./node_modules/.bin/tailwindcss build app.css -o public/styles/app.css

certs: generate_cert.go
	go run generate_cert.go --rsa-bits=2048 --host=localhost

clean:
	rm -rf data

cover:
	go tool cover -html=cover.out

down:
	docker-compose -p ${NAME} down

generate_cert.go:
	wget "https://raw.githubusercontent.com/golang/go/master/src/crypto/tls/generate_cert.go"

lint:
	golangci-lint run

migrate-create:
	migrate create -ext sql -dir storage/migrations -seq $(name)

migrate-down:
	migrate -path storage/migrations -database ${MIGRATE_DB_URL} down

migrate-goto:
	migrate -path storage/migrations -database ${MIGRATE_DB_URL} goto $(version)

migrate-up:
	migrate -path storage/migrations -database ${MIGRATE_DB_URL} up

start:
	go run cmd/server/*.go start

test:
	go test -coverprofile=cover.out -mod=readonly ./...

test-down:
	docker-compose -p ${NAME}-test -f docker-compose-test.yaml down

test-integration: test-up
	go test -p 1 -coverprofile=cover.out -tags=integration -mod=readonly ./...

test-up:
	docker-compose -p ${NAME}-test -f docker-compose-test.yaml up -d

up:
	mkdir -p data
	docker-compose -p ${NAME} up -d

up-build:
	mkdir -p data
	docker-compose -p ${NAME} up --build -d
