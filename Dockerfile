FROM golang:1.19-alpine as builder

WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o backup cmd/backup/main.go

FROM alpine:latest
RUN apk add --no-cache libc6-compat
COPY --from=builder /build/backup /bin/backup
COPY config_example.yaml /etc/backup/config.yaml

ENTRYPOINT ["backup"]
CMD ["schedule", "start"]