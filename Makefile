##
# Important! Before running any command make sure you have setup GOPATH:
# export GOPATH="$HOME/go"
# PATH="$GOPATH/bin:$PATH"

include Makefile.local

run:
	docker-compose up

stop:
	docker-compose down

integrationtest-run:
	go test ./integrationtest -v

unittest:
	go test $$(go list ./... | grep -v '/integrationtest') -v

# Creates new migration file with the current timestamp
# Example: make create-new-migration-file NAME=<name>
create-new-migration-file:
	$(eval NAME ?= noname)
	goose -dir ./schema/postgresql/migrations/ create $(NAME) sql

# go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.23
generate-sqlc:
	sqlc generate

migrate-up:
	goose -dir ./schema/postgresql/migrations -table schema_migrations postgres $(POSTGRES_DSN) up
migrate-redo:
	goose -dir ./schema/postgresql/migrations -table schema_migrations postgres $(POSTGRES_DSN) redo
migrate-down:
	goose -dir ./schema/postgresql/migrations -table schema_migrations postgres $(POSTGRES_DSN) down
migrate-reset:
	goose -dir ./schema/postgresql/migrations -table schema_migrations postgres $(POSTGRES_DSN) reset
migrate-status:
	goose -dir ./schema/postgresql/migrations -table schema_migrations postgres $(POSTGRES_DSN) status

rebuild: stop
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/api/main.go && docker-compose up -d --remove-orphans --build server

testmocks:
	# Generate mocks stubs.
	mockgen \
	-build_flags=--mod=mod \
	-destination=internal/mocks/mock_account_repository.go \
	-package mocks entain/internal/domain AccountRepository

	mockgen \
	-build_flags=--mod=mod \
	-destination=internal/mocks/mock_transaction_repository.go \
	-package mocks entain/internal/domain TransactionRepository

	mockgen \
	-build_flags=--mod=mod \
	-destination=internal/mocks/mock_tx_repository.go \
	-package mocks entain/internal/domain TxRepository

.DEFAULT_GOAL := run
