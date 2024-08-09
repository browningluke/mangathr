## Build
FROM golang:1.22-alpine AS build

WORKDIR /app

# Add build dependencies
RUN apk add make git

COPY ../../ .
RUN make

ENTRYPOINT ["mangathr"]

## Production
FROM alpine:latest

WORKDIR /

RUN mkdir config data

COPY --from=build /app/bin/mangathr /usr/local/bin/mangathr

ENTRYPOINT ["mangathr"]

