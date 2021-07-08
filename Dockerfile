FROM golang:alpine as build

LABEL maintainer "github.com/jxsl13"

WORKDIR /build
COPY *.go ./
COPY go.* ./

ENV CGO_ENABLED=0
ENV GOOS=linux 

RUN go get -d && go build -a -ldflags '-extldflags "-static"' -o discord-moderation .


FROM alpine:latest as minimal

# in container definitions
# set this to a path in order to use a .env file
ENV ENV_FILE ""

ENV BROKER_ADDRESS "rabbitmq:5672"
ENV BROKER_USER ""
ENV BROKER_PASSWORD ""
ENV REDIS_ADDRESS "redis:6379"
ENV REDIS_PASSWORD ""
ENV DATA_PATH "/data"
ENV BAN_REASON "VPN"
ENV BAN_DIRATION "24h"
ENV BROADCAST_BANS "false"
ENV BAN_COMMAND "ban {IP} {DURATION:MINUTES} {REASON}"

ENV DISCORD_TOKEN ""
ENV ADDRESS_CHANNEL_MAPPING ""
ENV LOGS_SKIP_JOIN_LEAVE "true"
ENV LOGS_SKIP_WHISPER "true"


WORKDIR /app
COPY --from=build /build/discord-moderation .
VOLUME ["/data"]
ENTRYPOINT ["/app/discord-moderation"]