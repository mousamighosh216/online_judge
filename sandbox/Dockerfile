FROM golang:1.21

WORKDIR /sandbox

# Install runtimes & compilers
RUN apt-get update && apt-get install -y \
    gcc \
    g++ \
    python3 \
    python3-pip \
    && rm -rf /var/lib/apt/lists/*

# Copy source
COPY . .

# Build executor worker
RUN go mod download
RUN go build -o executor ./cmd/internal/executor

# Runtime
CMD ["./executor"]
