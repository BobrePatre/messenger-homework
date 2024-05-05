FROM golang:latest AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server-service .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/server-service .

EXPOSE 8080

CMD ["./server-service"]
