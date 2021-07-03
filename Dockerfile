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
ENV BROKER_ADDRESS=rabbitmq:5672
ENV BROKER_USER=tw-admin
ENV BROKER_PASSWORD=""


WORKDIR /app
COPY --from=build /build/discord-moderation .
VOLUME ["/data"]
ENTRYPOINT ["/app/discord-moderation"]