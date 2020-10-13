FROM golang:1.15-buster AS builder
WORKDIR /src

# Copy go.mod and go.sum first, to enable caching of dependencies even after local code changes
COPY go.mod go.sum ./
RUN go mod download -x

COPY . ./
RUN go build -v -o /bin/server cmd/server/*.go

FROM debian:buster-slim
RUN set -x && apt-get update && \
  DEBIAN_FRONTEND=noninteractive apt-get install -y ca-certificates && \
  rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY development.toml ./
COPY --from=builder /bin/server ./

CMD ["./server", "-config", "development.toml"]
