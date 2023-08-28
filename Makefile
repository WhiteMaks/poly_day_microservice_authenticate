autogenerate-code:
	docker run --rm -v $(shell pwd):/src -w /src kjconroy/sqlc generate

test:
	docker-compose up -d
	timeout 60
	go clean -testcache && go test -cover -v ./...