FROM golang:latest AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o notification-service .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/notification-service .

EXPOSE 8080

CMD ["./notification-service"]
