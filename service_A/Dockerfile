# Dockerfile dentro de service_A/
# OBS: Se quiser 1.21, use golang:1.21
FROM golang:1.23

WORKDIR /app

# Copia os arquivos de mod
COPY go.mod go.sum ./
RUN go mod download

# Copia o restante
COPY . .

# Compila
RUN go build -o serviceA main.go

EXPOSE 8010

CMD ["./serviceA"]
