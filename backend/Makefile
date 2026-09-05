.PHONY: run build test migrate-up docker-up docker-down

run:
	go run cmd/server/main.go

build:
	go build -o bin/server cmd/server/main.go

test:
	go test ./...

migrate-up:
	# psql -d cafe_erp -f migrations/000001_init_schema.sql
	# psql -d cafe_erp -f migrations/000002_seed_data.sql

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down
