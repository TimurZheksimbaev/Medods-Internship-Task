FROM golang:1.22-alpine3.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .


RUN go build -o authentication_service ./main.go

FROM alpine:latest

WORKDIR /root/


COPY --from=builder /app/authentication_service .



RUN chmod +x ./authentication_service



CMD ["./authentication_service"]
