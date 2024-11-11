# Stage 1: Application Build
FROM golang:1.23.3-alpine3.20 AS builder

RUN apk update && apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download
COPY . .

RUN go build -o go-k8s-cache-demo ./cmd/go-k8s-cache-demo

# ---------------------------------------

# Stage 2: Docker image Build
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/go-k8s-cache-demo .

EXPOSE 8080

CMD ["./go-k8s-cache-demo"]