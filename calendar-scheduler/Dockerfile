# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Builder
FROM golang:1.13 as builder

RUN mkdir -p /opt/calendar
WORKDIR /opt/calendar

COPY ./script/wait-for-it.sh ./wait-for-it.sh
COPY ./config/config.docker.json ./config.json
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /opt/calendar/scheduler .

# Release
FROM bash:latest
LABEL maintainer="Grigoriy Zavodov <zavodov@gmail.com>"

RUN mkdir -p /calendar
WORKDIR /calendar

COPY --from=builder /opt/calendar/scheduler ./scheduler
COPY --from=builder /opt/calendar/config.json ./config.json
COPY --from=builder /opt/calendar/wait-for-it.sh ./wait-for-it.sh
RUN chmod +x ./wait-for-it.sh

CMD ["./wait-for-it.sh", "rabbit:15672", "--", "./scheduler", "--config=config.json"]
