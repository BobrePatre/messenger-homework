FROM golang:1.22.3 AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o service ./cmd/server

FROM scratch
# Alpine - dev only for docker health check, for prod use scratch
WORKDIR /root/
COPY --from=builder /app/service .

EXPOSE 8080

CMD ["./service"]
