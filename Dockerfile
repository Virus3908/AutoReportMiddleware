FROM golang:1.23.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN apt-get update && apt-get install -y unzip curl && \
    curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v24.3/protoc-24.3-linux-x86_64.zip && \
    unzip protoc-24.3-linux-x86_64.zip -d /usr/local && \
    rm -f protoc-24.3-linux-x86_64.zip

RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

RUN make generate

RUN make build-docker

FROM gcr.io/distroless/base-debian11 AS runner

WORKDIR /app

COPY --from=builder /app/bin/middleware /app/middleware

COPY config/config.yaml /app/config/config.yaml

ENV CONFIG_PATH=/app/config/config.yaml


ENTRYPOINT ["/app/middleware"]