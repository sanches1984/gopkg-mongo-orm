env:
	@cp .env.example .env
	@cp .env.example ./migrate/.env
	@docker run --name gopkg-test-mdb -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=password -p 5555:27017 -d mongo

test:
	go test -cover -v ./...