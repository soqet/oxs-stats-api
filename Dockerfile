FROM golang:alpine AS builder

ARG folder

ENV CGO_ENABLED=1

RUN apk update && apk upgrade
RUN apk add make
RUN apk add gcc \
    && apk add musl-dev

WORKDIR /app

COPY . .

WORKDIR /app/${folder}

RUN go build -o ./build/main.out ./cmd

FROM alpine:latest

EXPOSE $PORT

ARG folder

WORKDIR /app

COPY --from=builder ./app/${folder}/build/main.out .

CMD [ "./main.out" ]