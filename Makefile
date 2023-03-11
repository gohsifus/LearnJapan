run_dev:
	go run ./cmd/main.go

compose-up-mysql:
	docker-compose -f docker-compose.yml down -v --remove-orphans
	docker-compose up --build -d mysql

compose-up-pg:
	docker-compose -f docker-compose.yml down -v --remove-orphans
	docker-compose up --build -d pg
