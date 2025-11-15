FROM golang:1.25.1-alpine3.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY pkg/testutils ./pkg/testutils 
COPY tests ./tests

CMD ["go", "test", "-v", "./tests"]
