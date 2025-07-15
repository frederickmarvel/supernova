FROM golang:1.21 AS builder
WORKDIR /app
COPY go.mod go.sum .  
RUN go mod download  
COPY . .
RUN go build -o service ./cmd/service

FROM gcr.io/distroless/base
COPY --from=builder /app/service /service
ENTRYPOINT ["/service"]
