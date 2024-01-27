# build environment ###########################################
FROM golang:1.21.6-alpine@sha256:a49e5101be836613e432f0911b66a47dd48b50ebcd720717f70a30968237789e AS build-env

WORKDIR /app

# entrypoint
RUN apk add --no-cache entr
COPY entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]

# dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# server
COPY main.go .
COPY mutatingwebhook.go .
RUN go install .

# production image ############################################
FROM alpine:3.19.0@sha256:51b67269f354137895d43f3b3d810bfacd3945438e94dc5ac55fdac340352f48 AS prod
RUN apk add --no-cache ca-certificates

COPY --from=build-env /go/bin/network-health-admission-controller /bin/admission-controller

ENTRYPOINT ["admission-controller"]
