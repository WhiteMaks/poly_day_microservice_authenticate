install-dependencies:
	go install github.com/golang/mock/mockgen@v1.6.0

autogenerate-code:
	docker run --rm -v $(shell pwd):/src -w /src kjconroy/sqlc generate

autogenerate-mock:
	mockgen --package mockdatabase --destination autogenerate/database/mock/store.go poly_day_microservice_authenticate/autogenerate/database Store

autogenerate: autogenerate-code autogenerate-mock

up-test-environment:
	docker-compose up -d
	sleep 60

migrate-db:
	docker run --rm -v $(shell pwd)/database/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database postgres://user:password@localhost:15432/authenticate?sslmode=disable up

test:
	go clean -testcache && go test -cover -v ./...

build:
	go build -o ./.bin/authenticate/poly_day_microservice_authenticate cmd/main.go

setup-git:
	git config --global url."ssh://git@github.com".insteadOf "https://github.com"

.PHONY: build test migrate-db up-test-environment autogenerate-mock autogenerate-code install-dependencies setup-git