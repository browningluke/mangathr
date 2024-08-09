## Build
FROM golang:1.22-alpine AS build

WORKDIR /app

# Add build dependencies
RUN apk add make git

# Install go dependencies separately
# Allows for docker layer caching
COPY ../../go.mod ./
COPY ../../go.sum ./
COPY ../../Makefile ./
RUN make deps

# Build the package
COPY ../../ .
RUN make

ENTRYPOINT ["mangathr"]

## Production
FROM alpine:latest

WORKDIR /

RUN mkdir config data

COPY --from=build /app/bin/mangathr /usr/local/bin/mangathr

ENTRYPOINT ["mangathr"]

