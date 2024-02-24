FROM golang:1.21-alpine as builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod tidy

WORKDIR /app/cmd/app

COPY . .

RUN go build -o /app/main .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/.env .
COPY --from=builder /app/migrations ./migrations

ENTRYPOINT ["./main"]