# build environment ###########################################
FROM golang:1.19.0-alpine@sha256:7abc53c7b712049462918f8755d74739bf43bbf446a68789714ec0de2a506f18 AS build-env

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
FROM alpine:3.16.2@sha256:bc41182d7ef5ffc53a40b044e725193bc10142a1243f395ee852a8d9730fc2ad AS prod
RUN apk add --no-cache ca-certificates

COPY --from=build-env /go/bin/network-health-admission-controller /bin/admission-controller

ENTRYPOINT ["admission-controller"]
