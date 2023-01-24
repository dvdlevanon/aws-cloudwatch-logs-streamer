FROM golang:1.20-rc-bullseye AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
ADD pkg ./pkg
ADD cmd ./cmd

RUN go build -o /aws-cloudwatch-logs-streamer

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /aws-cloudwatch-logs-streamer /aws-cloudwatch-logs-streamer

USER nonroot:nonroot

ENTRYPOINT ["/aws-cloudwatch-logs-streamer"]
