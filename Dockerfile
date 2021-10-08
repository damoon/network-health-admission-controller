# build environment ###########################################
FROM golang:1.17.2-alpine@sha256:c2335038e2230960f81cb2f9f1fc5eca45e23b765de1848c7bbfaebcfd32d90d AS build-env

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
FROM alpine:3.14.2@sha256:e1c082e3d3c45cccac829840a25941e679c25d438cc8412c2fa221cf1a824e6a AS prod
RUN apk add --no-cache ca-certificates

COPY --from=build-env /go/bin/network-health-admission-controller /bin/admission-controller

ENTRYPOINT ["admission-controller"]
