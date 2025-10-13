FROM golang:latest AS development

ENV PATH="/go/bin:$PATH"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o server .

# Final minimal image
FROM alpine:latest
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=development /app/server .
COPY static ./static

# Expose the app port
EXPOSE 8080

# Environment variable for Mongo connection
ENV MONGO_URI=mongodb://mongo:27017/mydb

# Start the server
CMD ["./server"]