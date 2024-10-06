FROM golang:1.23 AS builder
RUN mkdir app
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o order-service .


FROM alpine:latest
WORKDIR /root/
RUN apk add libc6-compat
COPY --from=builder /app/order-service .
CMD ["./order-service"]