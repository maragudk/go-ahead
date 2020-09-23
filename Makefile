.PHONY: bindata build build-css clean cockroach-certs cockroach-sql cover down lint migrate-create migrate-goto migrate-up \
	start test test-down test-integration test-up up

export NAME := ahead
export VERSION := `git rev-parse --short HEAD`
export MIGRATE_DB_URL := "cockroachdb://root@localhost:26257/${NAME}?sslmode=verify-full&sslcert=certs/client.root.crt&sslkey=certs/client.root.key&sslrootcert=certs/ca.crt"

bindata:
	go-bindata -nometadata -pkg storage -o storage/migrations.go -ignore '.*\.go' -ignore '.DS_Store' -prefix storage/migrations storage/migrations

build:
	sed -i.bak "s/VERSION/${VERSION}/g" cmd/server/version.go
	GOOS=linux GOARCH=amd64 go build -o ${NAME} cmd/server/*.go
	sed -i.bak "s/${VERSION}/VERSION/g" cmd/server/version.go
	rm cmd/server/version.go.bak

build-css:
	NODE_ENV=production ./node_modules/.bin/tailwindcss build app.css -o public/styles/app.css

clean:
	rm -rf certs
	rm -rf cockroach-data

cockroach-certs:
	rm -rf certs
	mkdir -p certs/my-safe-directory
	cockroach cert create-ca --certs-dir=certs --ca-key=certs/my-safe-directory/ca.key
	cockroach cert create-node db localhost 127.0.0.1 --certs-dir=certs --ca-key=certs/my-safe-directory/ca.key
	cockroach cert create-client root --certs-dir=certs --ca-key=certs/my-safe-directory/ca.key
	cockroach cert create-client ${NAME} --certs-dir=certs --ca-key=certs/my-safe-directory/ca.key

cockroach-sql:
	cockroach sql --certs-dir certs

cover:
	go tool cover -html=cover.out

down:
	docker-compose -p ${NAME} down

lint:
	golangci-lint run

migrate-create:
	migrate create -ext sql -dir storage/migrations -seq $(name)

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
	mkdir -p cockroach-data
	docker-compose -p ${NAME} up -d
	sleep 1
	cockroach sql --certs-dir certs -e "create database if not exists ${NAME};"
	cockroach sql --certs-dir certs -e "create user if not exists ${NAME};"
	cockroach sql --certs-dir certs -e "grant select, insert, update, delete on database ${NAME} to ${NAME};"
