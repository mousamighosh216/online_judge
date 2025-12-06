FROM golang:1.21
WORKDIR /app

# Install compilers and Python
RUN apt-get update && \
    apt-get install -y gcc g++ python3 python3-pip && \
    rm -rf /var/lib/apt/lists/*

# Copy your executor project
COPY . .

# Download Go dependencies
RUN go mod download

# Build executor binary
RUN go build -o executor ./cmd/executor-worker/main.go

# Default command (optional)
CMD ["./executor"]
