FROM golang:1.17.1 as builder
WORKDIR /app
COPY go.mod go.sum /app/
RUN go mod download
COPY . /app/
RUN CGO_ENABLED=0 go build -o socks-proxy ./cmd/socks-proxy

FROM scratch
WORKDIR /app
COPY --from=builder /app/socks-proxy /app/
ENTRYPOINT ["/app/socks-proxy"]
