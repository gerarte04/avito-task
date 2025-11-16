FROM golang:1.25.1-alpine3.22 AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main cmd/main.go

FROM alpine:3.22 AS runner

RUN apk add --no-cache curl

WORKDIR /app

COPY --from=builder /build/main /app/main
COPY ./config/config.yaml /app/config.yaml
COPY ./docs/swagger-ui-dist /app/swagger-ui-dist

EXPOSE 8080

CMD ["/app/main", "--config=/app/config.yaml"]
