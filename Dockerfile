# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS builder
WORKDIR /app

COPY go.mod ./
COPY main.go ./
COPY templates ./templates
COPY static ./static
COPY imagenes ./imagenes

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server .

FROM alpine:3.22
WORKDIR /app

COPY --from=builder /app/server /app/server
COPY --from=builder /app/templates /app/templates
COPY --from=builder /app/static /app/static
COPY --from=builder /app/imagenes /app/imagenes

EXPOSE 8000
CMD ["/app/server", "-port=8000", "-dir=imagenes"]
