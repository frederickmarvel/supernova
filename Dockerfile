FROM golang:1.21 AS builder
WORKDIR /app
COPY go.mod go.sum .  
RUN go mod download  
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o service ./cmd/service

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/service /service
COPY --from=builder /app/.env /.env
WORKDIR /
RUN chmod +x /service
ENTRYPOINT ["/service"]
