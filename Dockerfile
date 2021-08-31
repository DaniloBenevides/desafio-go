#BUILD BINARY
FROM golang:1.17.0-alpine3.14 as builder


RUN mkdir build
COPY . /build
WORKDIR /build

RUN go build -o desafio

#BUILD IMAGE
FROM alpine:3.14.2

RUN mkdir -p /app
WORKDIR /app
COPY --from=builder build/desafio .
EXPOSE 8080
CMD ["./desafio"]
