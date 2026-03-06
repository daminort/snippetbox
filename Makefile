.PHONY: run
run:
	go mod tidy
	go fmt ./...
	go run ./cmd/web

.PHONY: test
test:
	go mod tidy
	go fmt ./...
	go test ./cmd/web

.PHONY: db-start
db-start:
	docker-compose up -d

.PHONY: db-stop
db-stop:
	docker-compose down
