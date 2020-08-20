.PHONY: bindata build clean cockroach-certs cockroach-sql cover down lint migrate-create migrate-up start test test-integration up

export VERSION := `git rev-parse --short HEAD`
export MIGRATE_DB_URL := "cockroachdb://root@localhost:26257/ahead?sslmode=verify-full&sslcert=certs/client.root.crt&sslkey=certs/client.root.key&sslrootcert=certs/ca.crt"

bindata:
	go-bindata -nometadata -pkg storage -o storage/migrations.go -ignore '.*\.go' -ignore '.DS_Store' -prefix storage/migrations storage/migrations

build:
	sed -i.bak "s/VERSION/${VERSION}/g" cmd/ahead/ahead.go
	GOOS=linux GOARCH=amd64 go build -o ahead cmd/ahead/*.go
	sed -i.bak "s/${VERSION}/VERSION/g" cmd/ahead/ahead.go
	rm cmd/ahead/ahead.go.bak

clean:
	rm -rf certs
	rm -rf cockroach-data

cockroach-certs:
	rm -rf certs
	mkdir -p certs/my-safe-directory
	cockroach cert create-ca --certs-dir=certs --ca-key=certs/my-safe-directory/ca.key
	cockroach cert create-node db localhost 127.0.0.1 --certs-dir=certs --ca-key=certs/my-safe-directory/ca.key
	cockroach cert create-client root --certs-dir=certs --ca-key=certs/my-safe-directory/ca.key
	cockroach cert create-client ahead --certs-dir=certs --ca-key=certs/my-safe-directory/ca.key

cockroach-sql:
	cockroach sql --certs-dir certs

cover:
	go tool cover -html=cover.out

down:
	docker-compose -p ahead down

lint:
	golangci-lint run

migrate-create:
	migrate create -ext sql -dir storage/migrations -seq $(name)

migrate-goto:
	migrate -path storage/migrations -database ${MIGRATE_DB_URL} goto $(version)

migrate-up:
	migrate -path storage/migrations -database ${MIGRATE_DB_URL} up

start: up
	go run cmd/ahead/*.go start

test:
	go test -coverprofile=cover.out -mod=readonly ./...

test-integration:
	go test -p 1 -coverprofile=cover.out -tags=integration -mod=readonly ./...

up:
	mkdir -p cockroach-data
	docker-compose -p ahead up -d
	cockroach sql --certs-dir certs -e 'create database if not exists ahead;'
	cockroach sql --certs-dir certs -e 'create user if not exists ahead;'
	cockroach sql --certs-dir certs -e 'grant select, insert, update, delete on database ahead to ahead;'
