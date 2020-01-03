# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Builder
FROM golang:1.13 as builder

RUN mkdir -p /opt/calendar
RUN mkdir -p /opt/script
WORKDIR /opt/calendar

COPY ./script/wait-for-it.sh /opt/script/wait-for-it.sh
COPY ./config/config.docker.json ./config.json
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /opt/calendar/api .

# Release
FROM alpine:latest
LABEL maintainer="Grigoriy Zavodov <zavodov@gmail.com>"

RUN mkdir -p /calendar
WORKDIR /calendar

COPY --from=builder /opt/calendar ./
COPY --from=builder /opt/script/wait-for-it.sh /usr/bin/wait-for-it.sh
RUN chmod +x /usr/bin/wait-for-it.sh
RUN apk update && apk add bash && apk add --no-cache coreutils

EXPOSE 8888
ENTRYPOINT ["./api", "--config=config.json"]
