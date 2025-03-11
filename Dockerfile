FROM golang:1.23.4

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /feederizer

# Copy Go module files and download dependencies
COPY server/go.mod server/go.sum .server/
RUN go mod download # TODO

# Copy the rest of the application code
# TODO:
COPY . .

# Build the Go application
# TODO:
RUN go build -o server ./cmd/main.go

# Expose the port the application runs on
EXPOSE 8080

# Run the application
# TODO:
CMD ["./server"]
