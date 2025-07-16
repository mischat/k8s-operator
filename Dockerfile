FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY programA.go .
RUN go build -o programA programA.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/programA .
CMD ["./programA"]
