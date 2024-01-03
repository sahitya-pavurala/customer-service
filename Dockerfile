# Dockerfile
FROM golang:1.21-alpine as builder

WORKDIR /app
COPY . .

RUN go build -o customer-service .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/customer-service /app/

CMD ["./customer-service"]
