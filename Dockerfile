FROM golang:1.17-alpine

WORKDIR /app

COPY . .

RUN go mod download

CMD ["go", "run", "main.go"]
