## Build
FROM golang:1.23-alpine AS build

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

## Production
FROM gcr.io/distroless/static-debian12

WORKDIR /
COPY --from=build /app/bin/mangathr /mangathr

VOLUME ["/config", "/data"]
ENTRYPOINT ["/mangathr"]
