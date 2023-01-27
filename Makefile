run_dev:
	go run ./cmd/main.go

compose-up-db:
	docker-compose -f docker-compose.yml down -v --remove-orphans
	docker-compose up --build -d db
