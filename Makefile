include ./config/development.env

OUT_DIR=$(PWD)/pkg/pb
POSTGRESQL_URL="postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable"
MIGRATE=migrate -database ${POSTGRESQL_URL} -path ./migrations
PROTO=docker run --rm -v "$(PWD):/defs" namely/protoc-all -f ./pkg/pb/auth.proto -o ./ -l go

.SILENT: daemon mockgen protogen lint compose test
.PHONY: migrate-up migrate-down migrate-new migrate-force

daemon:
	APP_ENV=development CompileDaemon \
		--build="go build -o ./bin/auth ./cmd/auth/main.go" \
		--command="./bin/auth" \
		-graceful-kill \
		-log-prefix=false \
		-polling \
		-polling-interval=350

migrate-up:
	@echo "Running all new database migrations..."
	@$(MIGRATE) up

migrate-new:
	@read -p "Enter the name of new migration: " name; \
	$(MIGRATE) create -ext sql -dir migrations $$name

migrate-down:
	@echo "Running all down database migrations..."
	@$(MIGRATE) down

migrate-force:
	@read -p "Enter version migration: " name; \
	$(MIGRATE) force $$name

migrate-drop:
	@echo "Dropping everything in database..."
	@$(MIGRATE) drop

codegen:
	go generate -v ./...

lint:
	golangci-lint run -c .golangcilint.yaml

compose:
	docker compose -f ./deployments/docker-compose.yml -p auth up

protogen:
	@echo "Running generation protofiles..."
	@$(PROTO)

test:
	mkdir -p ./tmp
	go clean -testcache
	go test -v ./... -coverprofile=./tmp/coverage.out

coverage:
	gotestsum --format dots --packages="./..." -- -coverprofile=./tmp/coverage.out
	go tool cover -html=./tmp/coverage.out -o ./tmp/coverage.htmlgateway

build:
	GOARCH=amd64 \
	GOOS=linux \
	CGO_ENABLED=0 \
	APP_ENV=development \
	go build -o ./bin/auth ./cmd/auth/main.go