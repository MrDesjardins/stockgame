# PHONY does not check if the file exists or not, it will always run the command
.PHONY: api 
.PHONY: init 
.PHONY: unit-test 

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
	go test ./... -coverpkg=./... -coverprofile=./coverage/coverage.out
	go tool cover -func ./coverage/coverage.out
# go tool cover -html=coverage.out -o coverage.html

db:
	duckdb data/db/stockgame.duckdb