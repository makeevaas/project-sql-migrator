# -----------------
# commands
# -----------------
.PHONY: help
help: Makefile
	@echo "Choose a command in:"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'

## build - сборка приложения
.PHONY: build
build:
	@echo "Сборка утилиты  sql-migrator"
	@DB_CONNECTION_PATH="postgres://postgres:pwd@localhost:5432/main_db?sslmode=disable"
	@MIGRATIONS_PATH=./migrations
	@go build -o ./bin/migrate cmd/main.go

## postgres - скачивание и запуск postgres
.PHONY: postgres
postgres:
	@echo "Cкачивание и запуск postgres:17.0"
	@docker pull postgres:17.0
	@docker compose up
   
## migrate_create - создание файла миграции
.PHONY: migrate_create
migrate_create:
	@echo "Выполняется команда migrate_create"
	@./bin/migrate -create

## migrate_up - накат миграции
.PHONY: migrate_up
migrate_up:
	@echo "Выполняется команда migrate_up"
	@./bin/migrate -up

## migrate_down - откат миграции
.PHONY: migrate_down
migrate_down:
	@echo "Выполняется команда migrate_down"
	@./bin/migrate -down

## migrate_redo - повтор последней миграции
.PHONY: migrate_redo
migrate_redo:
	@echo "Выполняется команда migrate_redo"
	@./bin/migrate -redo

## migrate_status - получение статуса миграций
.PHONY: migrate_status
migrate_status:
	@echo "Выполняется команда migrate_status"
	@./bin/migrate -status		

## test - запуск тестов
.PHONY: test
test:
	@echo "Запуск тестов"
	@go test -v	./internal/tests