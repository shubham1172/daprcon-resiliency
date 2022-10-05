## Build
FROM golang:1.19-alpine3.16 AS build

WORKDIR /app

COPY go.* ./
COPY *.go ./

RUN go build -o /availability

## Deploy
FROM alpine:3.16

WORKDIR /

COPY --from=build /availability /availability

EXPOSE 8000

ENTRYPOINT ["/availability"]
