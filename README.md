# BLOG BACKEND SERVER

Simple REST API project using Go with Clean Architecture

## What has been used:
* [echo](https://github.com/labstack/echo) - Web framework
* [sqlx](https://github.com/jmoiron/sqlx) - Extensions to database/sql.
* [pgx](https://github.com/jackc/pgx) - PostgreSQL driver and toolkit for Go
* [viper](https://github.com/spf13/viper) - Go configuration with fangs
* [go-redis](https://github.com/go-redis/redis) - Type-safe Redis client for Golang
* [zap](https://github.com/uber-go/zap) - Logger
* [validator](https://github.com/go-playground/validator) - Go Struct and Field validation
* [paseto](https://github.com/o1egl/paseto) - Platform-Agnostic Security Tokens (PASETO)
* [uuid](https://github.com/google/uuid) - UUID
* [migrate](https://github.com/golang-migrate/migrate) - Database migrations. CLI and Golang library.
* [minio-go](https://github.com/minio/minio-go) - MinIO Client SDK for Go
* [swag](https://github.com/swaggo/swag) - Swagger
* [testify](https://github.com/stretchr/testify) - Testing toolkit
* [gomock](https://github.com/golang/mock) - Mocking framework

## Quick start

Clone this repository:
```sh
git clone https://github.com/scul0405/blog-clean-architecture-rest-api.git
cd blog-clean-architecture-rest-api
```
#### Docker development usage

Run this command:
```sh
make docker_dev
```
#### Local development usage

```sh
make docker_local
```

## Local development

### Install tools

- [Golang](https://golang.org/)
- [Docker desktop](https://www.docker.com/products/docker-desktop)
- [TablePlus](https://tableplus.com/)
- [Migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
- [Gomock](https://github.com/golang/mock)

Run server
```sh
make run
```

Run testing
```sh
make test
```

Generate documentation
```sh
make swag
```

## Documentation

### Swagger UI
[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

## Monitor

### Jaeger
[http://localhost:16686](http://localhost:16686)

### Minio
[http://localhost:9001](http://localhost:9001)

