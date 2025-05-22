FROM golang:1.24.3 AS builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o volume-cleaner .

FROM alpine:3.21

COPY --from=builder /app/volume-cleaner /volume-cleaner
ENTRYPOINT ["/bin/sh"]
