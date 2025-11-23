FROM golang:1.24.3-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o review-assigner ./cmd/main.go

EXPOSE 8080

CMD ["./review-assigner"]