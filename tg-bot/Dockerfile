FROM golang:1.12 AS builder
WORKDIR /go/src/github.com/lokeon-university/kfr-telegram-bot
COPY . .
RUN go get ./...
RUN CGO_ENABLED=0 GOOS=linux go build -v -o kfr-bot
FROM golang:alpine AS certs
RUN apk --no-cache add ca-certificates
FROM alpine
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/lokeon-university/kfr-telegram-bot/kfr-bot /kfr-bot
CMD ["/kfr-bot"]