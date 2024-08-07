## Build
FROM golang:1.22-alpine AS build

WORKDIR /app

COPY ../../go.mod ./
COPY ../../go.sum ./
RUN go mod download && go mod verify

COPY ../../ .
RUN go build -v -o /mangathr ./cmd/mangathr

ENTRYPOINT ["mangathr"]

## Production
FROM alpine:latest

WORKDIR /

RUN mkdir config data

COPY --from=build /mangathr /mangathr

ENTRYPOINT ["/mangathr"]
