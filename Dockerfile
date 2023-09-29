FROM golang:1.21.1 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /app/pastebin ./cmd/web/main.go

EXPOSE 8080

CMD ["/app/pastebin"]
