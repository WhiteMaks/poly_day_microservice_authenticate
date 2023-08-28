install-dependencies:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

autogenerate:
	docker run --rm -v $(pwd):/src -w /src kjconroy/sqlc generate

test:
	go clean -testcache && go test -cover -v ./...