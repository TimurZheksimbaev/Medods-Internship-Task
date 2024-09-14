FROM golang:1.22-bookworm AS builder

ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o main .

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/main .

RUN chmod +x ./main

COPY app.env ./app.env

EXPOSE 3000
EXPOSE 5432

CMD ["./main"]
