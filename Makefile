autogenerate-code:
	docker run --rm -v $(shell pwd):/src -w /src kjconroy/sqlc generate

up-test-environment:
	docker-compose up -d
	sleep 60

migrate-db:
	docker run --rm -v $(shell pwd)/database/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database postgres://user:password@localhost:15432/authenticate?sslmode=disable up

test:
	go clean -testcache && go test -cover -v ./...