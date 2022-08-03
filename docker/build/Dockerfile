## Build
FROM golang:1.18-alpine AS build

WORKDIR /app

COPY ../../go.mod ./
COPY ../../go.sum ./
RUN go mod download && go mod verify

COPY ../../ .
RUN go build -v -o /mangathrv2 .

ENTRYPOINT ["mangathrv2"]

## Production
FROM alpine:latest

WORKDIR /

RUN mkdir config data

COPY --from=build /mangathrv2 /mangathrv2

ENTRYPOINT ["/mangathrv2"]
