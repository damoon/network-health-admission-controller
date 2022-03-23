# build environment ###########################################
FROM golang:1.18.0-alpine@sha256:9ccb0ed869157f3b1630f5dda3422b5974defa9dd82c7375ca68dc3a9cbf8fae AS build-env

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
FROM alpine:3.15.2@sha256:66b861b1099af1551a0eee163c175fd008744192c3fbb7f22e998db0ce09e8ea AS prod
RUN apk add --no-cache ca-certificates

COPY --from=build-env /go/bin/network-health-admission-controller /bin/admission-controller

ENTRYPOINT ["admission-controller"]
