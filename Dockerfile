FROM golang:1.24.2 AS builder

COPY . /src
WORKDIR /src

RUN go build -o server ./cmd/server/main.go

ENV SCALABLE_FLAKE_DB_ADDR=127.0.0.1:6379
ENV SCALABLE_FLAKE_BACKEND=redis
ENV SCALABLE_FLAKE_TENANT=default

EXPOSE 8000
EXPOSE 9000

CMD ["./server"]
