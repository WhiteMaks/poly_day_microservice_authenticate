docker run --rm -v ${PWD}:/src -w /src kjconroy/sqlc generate
Start-Sleep -Seconds 5
mockgen --package mockdatabase --destination autogenerate/database/mock/store.go poly_day_microservice_authenticate/autogenerate/database Store