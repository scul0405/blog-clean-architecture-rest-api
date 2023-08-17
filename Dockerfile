FROM golang:1.20-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o main cmd/api/main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder app/migrate ../bin/migrate
COPY config/config-docker.yaml ./config/config.yaml
COPY migrations ./migrations
COPY start.sh .
RUN chmod +x start.sh
COPY wait-for.sh .
RUN chmod +x wait-for.sh

EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]
