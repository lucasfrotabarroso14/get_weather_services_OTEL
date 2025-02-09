# Dockerfile dentro de service_B/
FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o serviceB main.go

EXPOSE 8091

CMD ["./serviceB"]

