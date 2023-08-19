# ==============================================================================
# Main

run:
	go run ./cmd/api/main.go

build:
	go build ./cmd/api/main.go

test:
	go test -cover ./...

# ==============================================================================
# golang-migrate postgresql

force:
	migrate -database postgres://admin:secret@localhost:5432/blog_db?sslmode=disable -path migrations force 1

version:
	migrate -database postgres://admin:secret@localhost:5432/blog_db?sslmode=disable -path migrations version

migrate_up:
	migrate -database postgres://admin:secret@localhost:5432/blog_db?sslmode=disable -path migrations up 1

migrate_down:
	migrate -database postgres://admin:secret@localhost:5432/blog_db?sslmode=disable -path migrations down 1

# ==============================================================================
# Modules support

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

deps-upgrade:
	go get -u -t -d -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache

# ==============================================================================
# Docker compose commands
docker_dev:
	docker-compose -f docker-compose.dev.yaml up --build

docker_local:
	docker-compose -f docker-compose.local.yaml up --build

# ==============================================================================