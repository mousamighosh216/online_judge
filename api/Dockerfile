# 1️⃣ Base image: Go installed on Linux
FROM golang:1.24

# 2️⃣ Set working directory inside container
WORKDIR /app

# 3️⃣ Copy go.mod and go.sum first (for caching)
COPY go.mod go.sum ./

# 4️⃣ Download dependencies
RUN go mod download

# 5️⃣ Copy the rest of the source code
COPY . .

# 6️⃣ Build the server binary
RUN go build -o server ./cmd/server

# 7️⃣ Expose API port
EXPOSE 8080

# 8️⃣ Run the server
CMD ["./server"]
