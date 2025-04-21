# PHONY does not check if the file exists or not, it will always run the command
.PHONY: watch 
.PHONY: api 
.PHONY: init 
.PHONY: unit-test 
.PHONY: unit-test-watch 
.PHONY: unit-test-coverage
.PHONY: db
.PHONY: web-release
.PHONY: go-release
.PHONY: sync-env
.PHONY: generate-constants

dev: sync-env generate-constants
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
	find ./data/raw/stocks/ -type f -name '*_cleaned*' -exec rm -f {} +
	go run cmd/data-loader/main.go

unit-test: 
#	go test ./... 
	gotestsum --format testname

unit-test-watch: 
	gotestsum --format testname --watch

unit-test-coverage:
	go test ./... -coverpkg=./... -coverprofile=./coverage/coverage.out || true
	go tool cover -func=./coverage/coverage.out || true
	go tool cover -html=./coverage/coverage.out -o=./coverage/coverage.html

db:
	docker-compose -f docker-compose.yml up -d

container-debug:
	docker exec -it stock_postgres bash 

db-debug:
	PGPASSWORD=mypassword psql -h localhost -p 5432 -U myuser -d mydb

web-release:
	@echo "Running web release build..."
	(cd cmd/frontend-server && . ~/.nvm/nvm.sh && nvm use && npm run build)
	@echo "Running web go server..."
	go build -o bin/api-server cmd/api-server/main.go
	@echo "Move static files to bin folder..."
	mkdir -p bin/assets
	cp -r cmd/frontend-server/dist/* bin
	cp .env bin/.env
	
go-release: generate-constants
	@echo "Running release build..."
	go build -o bin/api-server cmd/api-server/main.go
	go build -o bin/data-loader cmd/data-loader/main.go

sync-env:
	@echo "Running sync-env..."
	cp .env ./cmd/frontend-server/.env

generate-constants:
	@echo "Running generate-constants..."
	go run cmd/back-to-front/main.go