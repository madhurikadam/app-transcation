FROM golang:1.17-alpine

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /app

COPY .build/app-transcation  ./

ENTRYPOINT ["./app-transcation"]
