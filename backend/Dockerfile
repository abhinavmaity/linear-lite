FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/migrate ./cmd/migrate

FROM alpine:3.22

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /bin/api /usr/local/bin/api
COPY --from=builder /bin/migrate /usr/local/bin/migrate
COPY migrations ./migrations

EXPOSE 8080

CMD ["api"]

