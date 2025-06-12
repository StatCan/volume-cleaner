FROM golang:1.24.3 AS scheduler-builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o volume-cleaner-scheduler ./cmd/scheduler

FROM alpine:3.21

COPY --from=scheduler-builder /app/volume-cleaner-scheduler /volume-cleaner-scheduler
ENTRYPOINT ["/bin/sh"]
