install-dependencies:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

autogenerate:
	sqlc generate

test:
	go clean -testcache && go test -cover -v ./...