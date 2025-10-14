FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .
COPY static ./static

EXPOSE 8080

ENV MONGO_URI=mongodb://mongo:27017/mydb

CMD ["./server"]
