FROM golang:1.16

WORKDIR /app

# Copy the Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o grpc-server

EXPOSE 50051

CMD ["./grpc-server"]
