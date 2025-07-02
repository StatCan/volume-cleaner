FROM golang:1.24.3 AS controller-builder

WORKDIR /app
COPY ./volume-cleaner/ .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o volume-cleaner-controller ./cmd/controller

FROM alpine:3.21

COPY --from=controller-builder /app/volume-cleaner-controller /volume-cleaner-controller
ENTRYPOINT ["/bin/sh"]
