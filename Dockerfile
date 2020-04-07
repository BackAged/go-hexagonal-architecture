# Defining App builder image
FROM golang:alpine AS builder

ARG VERSION=unversioned
ARG PRIVATE_KEY

# Add git to determine build git version
RUN apk add --no-cache --update git openssh

# Set GOPATH to build Go app
ENV GOPATH=/go

# Set apps source directory
ENV SRC_DIR=${GOPATH}/src/github.com/BackAged/go-hexagonal-architecture

# Define current working directory
WORKDIR ${SRC_DIR}

# Copy apps scource code to the image
COPY . ${SRC_DIR}

# Build App
RUN ./build.sh

# Defining App image
FROM alpine:latest

RUN apk add --no-cache --update ca-certificates

# Copy App binary to image
COPY --from=builder /go/bin/go-hexa /usr/local/bin/go-hexa

EXPOSE 8000

ENTRYPOINT ["go-hexa"]
