# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Builder
FROM golang:1.13 as builder

RUN mkdir -p /opt/banner-rotation
RUN mkdir -p /opt/script
WORKDIR /opt/banner-rotation

COPY ./script/wait-for-it.sh /opt/script/wait-for-it.sh
COPY ./config/config.docker.json ./config.json
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /opt/banner-rotation/app .

# Release
FROM alpine:latest
LABEL maintainer="Grigoriy Zavodov <zavodov@gmail.com>"

RUN mkdir -p /banner-rotation
WORKDIR /banner-rotation

COPY --from=builder /opt/banner-rotation ./
COPY --from=builder /opt/script/wait-for-it.sh /usr/bin/wait-for-it.sh
RUN chmod +x /usr/bin/wait-for-it.sh
RUN apk update && apk add bash && apk add --no-cache coreutils

EXPOSE 8390
ENTRYPOINT ["./app", "--config=config.json"]
