FROM golang:1 AS builder

WORKDIR /app

COPY . .

RUN go build -o /blitz-proxy .

FROM scratch

COPY --from=builder /blitz-proxy /blitz-proxy

CMD ["/blitz-proxy"]
