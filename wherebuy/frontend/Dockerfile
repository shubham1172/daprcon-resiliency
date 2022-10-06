## Build
FROM golang:1.19-alpine3.16 AS build

WORKDIR /app

COPY go.* ./
COPY *.go ./

RUN GOOS=linux GOARCH=amd64 go build -o /frontend

## Deploy
FROM alpine:3.16

WORKDIR /

COPY --from=build /frontend /frontend

EXPOSE 9000

ENTRYPOINT ["/frontend"]
