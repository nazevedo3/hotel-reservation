build:
	@go build -o bin/api

run: build
	@./bin/api

seed:
	@go run scripts/seed.go

test:
	@go test -v ./...

mongo:
	docker run --name mongodb -d -p 27017:27017 mongo:latest