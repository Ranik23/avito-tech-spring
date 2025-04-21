FROM golang:1.24.1-alpine AS builder


WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o bin ./cmd/main/main.go


FROM alpine:latest


WORKDIR /app

COPY --from=builder /app/bin .

COPY ./migrations ./migrations

COPY ./config ./config
COPY .env .env


EXPOSE 8080

ENTRYPOINT ["./bin"]