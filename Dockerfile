FROM mcr.microsoft.com/devcontainers/go:0-1-bullseye

# Setup working dorectory
WORKDIR /app

# We are copying this before everything else because it's faster
COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
COPY . .

# Set environment variables
EXPOSE 8080

# Runtime operations (CMD)
CMD go run *.go