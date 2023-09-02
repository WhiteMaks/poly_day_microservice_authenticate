# PolyDay Microservice Authenticate

## Развертка проекта в докер контейнере

Создать файл docker-compose.yml:


```yaml
version: "3.8"

services:

  poly_day_microservice_authenticate_postgresql:
    image: postgres:14.2-alpine
    container_name: poly_day_microservice_authenticate_postgresql
    environment:
      POSTGRES_DB: "authenticate"
      POSTGRES_USER: "user"
      POSTGRES_PASSWORD: "password"
      PGDATA: "/var/lib/postgresql/data/.database"
    volumes:
      - .:/var/lib/postgresql/data
    ports:
      - "15432:5432"

  poly_day_microservice_authenticate_migrator:
    image: whitemaks/poly_day_microservice_authenticate_migrator:release_0.0.6
    container_name: poly_day_microservice_authenticate_migrator
    environment:
      POSTGRES_CONNECTING_STRING: "postgres://user:password@poly_day_microservice_authenticate_postgresql:5432/authenticate?sslmode=disable"
    restart: on-failure

  poly_day_microservice_authenticate_app:
    image: whitemaks/poly_day_microservice_authenticate:release_0.0.6
    container_name: poly_day_microservice_authenticate_app
    environment:
      POSTGRES_CONNECTING_STRING: "postgres://user:password@poly_day_microservice_authenticate_postgresql:5432/authenticate?sslmode=disable"
      APP_PORT: "8080"
    ports:
      - "8080:8080"
    restart: on-failure
```

и выполнить команду:
```shell
docker-compose up -d
```