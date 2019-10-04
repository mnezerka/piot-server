FROM golang:alpine AS builder

RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s"
CMD ["./piot-server"]

FROM alpine:latest AS alpine
COPY --from=builder /app/piot-server /app/piot-server
WORKDIR /app/
EXPOSE 8000
CMD ["./piot-server"]
