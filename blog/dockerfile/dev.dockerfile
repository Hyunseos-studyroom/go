FROM golang:1.23.2

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/githubnemo/CompileDaemon@latest

COPY . .

EXPOSE 8080

CMD ["CompileDaemon", "-build", "go build -o main .", "-command", "./main"]
