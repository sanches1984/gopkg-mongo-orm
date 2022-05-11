env:
	@cp .env.example .env
	@cp .env.example ./migrate/.env

container:
	@docker run --name gopkg-test-mdb -p 5555:27017 -d mongo

test:
	go test -cover -v --tags=ci ./...

test-all:
	go test -cover -v ./...