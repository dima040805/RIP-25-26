.PHONY: run migrate test vet swag up down certs

run:
	go run ./cmd/PlanetsRadiusResearch

migrate:
	go run ./cmd/migrate

test:
	go test -race ./...

vet:
	go vet ./...

swag:
	swag init -g cmd/PlanetsRadiusResearch/main.go -o docs

up:
	docker compose up --build -d

down:
	docker compose down

certs:
	openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout certs/server.key -out certs/server.crt -config certs/server.cnf
