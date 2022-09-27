# syntax=docker/dockerfile:1

## Build
FROM golang:1.18-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /bsc

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /bsc /bsc

EXPOSE 8000
EXPOSE 9000

USER nonroot:nonroot

ENTRYPOINT ["/bsc"]
