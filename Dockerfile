# Use the official Golang image as the builder
FROM golang:1.20 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy source files into the container
COPY . .

# Initialize Go module if go.mod does not exist
RUN [ ! -f go.mod ] && go mod init elasticsearch-inserter || echo "go.mod exists"
RUN [ ! -f go.sum ] && go mod tidy || echo "go.sum exists"

# Set CGO_ENABLED=0 to ensure compatibility with Alpine and build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o elasticsearch-inserter main.go

# Use a minimal image for running the binary
FROM alpine:latest

# Set the timezone (optional)
RUN apk add --no-cache tzdata

# Copy the binary from the builder stage
COPY --from=builder /app/elasticsearch-inserter /usr/local/bin/elasticsearch-inserter

# Ensure the binary is executable
RUN chmod +x /usr/local/bin/elasticsearch-inserter

# Set the entrypoint to the executable
ENTRYPOINT ["/usr/local/bin/elasticsearch-inserter"]

# Set default values for Elasticsearch flags (can be overridden at runtime)
CMD ["-es_host=http://localhost:9200", "-username=elastic", "-password=changeme", "-indexname=my-index", "-tool=tool"]

