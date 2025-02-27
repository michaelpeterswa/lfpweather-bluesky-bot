# -=-=-=-=-=-=- Compile Image -=-=-=-=-=-=-

FROM golang:1 AS stage-compile

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./... && CGO_ENABLED=0 GOOS=linux go build ./cmd/lfpweather-bluesky-bot

# -=-=-=-=- Final Distroless Image -=-=-=-=-

# hadolint ignore=DL3007
FROM gcr.io/distroless/static-debian12:latest AS stage-final

COPY --from=stage-compile /go/src/app/lfpweather-bluesky-bot /
CMD ["/lfpweather-bluesky-bot"]