FROM golang:1.22.6-alpine3.19 as builder

WORKDIR /app
COPY . .

RUN go mod download
RUN go mod verify

RUN go build -o main main.go

FROM debian as runner

RUN apt-get update && apt-get install -y \
  ca-certificates \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /app/main ./main

CMD ["/app/main"]
