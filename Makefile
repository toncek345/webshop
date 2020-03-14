

.PHONY: web
web:
	go run ./cmd/web/main.go

.PHONY: migrate
migrate:
	go run ./cmd/migrate/main.go
