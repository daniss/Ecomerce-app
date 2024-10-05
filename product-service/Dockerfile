FROM golang:1.23 AS builder
RUN mkdir app
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o product-service .


FROM alpine:latest
WORKDIR /root/
RUN apk add libc6-compat
COPY --from=builder /app/product-service .
CMD ["./product-service"]