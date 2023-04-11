FROM mcr.microsoft.com/devcontainers/go:0-1-bullseye

# Setup working dorectory
WORKDIR /app

# We are copying this before everything else because it's faster
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

# Set environment variables

# Runtime operations (CMD)
CMD go run *.go