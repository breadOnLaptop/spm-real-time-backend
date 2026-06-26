FROM golang:latest AS builder

# Install protoc
RUN apt-get update && apt-get install -y protobuf-compiler

# Install go protoc plugins
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Generate gRPC bindings
RUN protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/telemetry.proto

RUN CGO_ENABLED=0 go build -o spm-api ./cmd/server/main.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/spm-api .

EXPOSE 8000 50051
CMD ["./spm-api"]
