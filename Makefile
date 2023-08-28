autogenerate-code:
	docker run --rm -v $(shell pwd):/src -w /src kjconroy/sqlc generate

test:
	go clean -testcache && go test -cover -v ./...