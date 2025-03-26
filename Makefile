# PHONY does not check if the file exists or not, it will always run the command
.PHONY: watch 
.PHONY: api 
.PHONY: init 
.PHONY: unit-test 
.PHONY: unit-test-coverage
.PHONY: db
.PHONY: release
.PHONY: generate-constants

dev: generate-constants
	air -c .air.toml & \
	(cd cmd/frontend-server && . ~/.nvm/nvm.sh && nvm use && npm run dev) & \
	wait

api-watch:
	air -c .air.toml

api:
	go run cmd/api-server/main.go

init:
	go mod tidy
	go mod verify
	go build
	go run cmd/data-loader/main.go

unit-test: 
	go test -parallel 1 ./... 

unit-test-coverage:
	go test ./... -coverpkg=./... -coverprofile=./coverage/coverage.out || true
	go tool cover -func=./coverage/coverage.out || true
	go tool cover -html=./coverage/coverage.out -o=./coverage/coverage.html

db:
	duckdb data/db/stockgame.duckdb

release: generate-constants
	go build -o bin/api-server cmd/api-server/main.go
	go build -o bin/data-loader cmd/data-loader/main.go

generate-constants:
	go run cmd/back-to-front/main.go