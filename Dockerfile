# syntax = docker/dockerfile:1

########################################
## Build Stage
########################################
FROM golang:1.22-alpine3.19 as builder

ARG VERSION=v0.0.1
ENV VERSION=${VERSION}

# add a label to clean up later
LABEL stage=intermediate

# install required packages
RUN apk update && apk add --no-cache git tzdata gcc musl-dev vips-dev

# Install wire
RUN go install github.com/google/wire/cmd/wire@latest

# setup the working directory
WORKDIR /src

# COPY go.mod and go.sum files to the workspace
COPY go.mod .
COPY go.sum .

# Get dependencies - will also be cached if we won't change mod/sum
RUN go mod download

# add source code
ADD . .

# build the source
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X github.com/praveenmsp23/trackdocs/pkg/config.BuildSHA=local -X github.com/praveenmsp23/trackdocs/pkg/config.BuildBranch=local -X github.com/praveenmsp23/trackdocs/pkg/config.BuildTime=$(date +%s) " -o /out/trackdocs-api cmd/api/main.go cmd/api/inject_*.go cmd/api/wire_gen.go

########################################
## Production Stage for trackdocs-api
########################################
FROM alpine:3.19 as trackdocs-api

RUN apk update && apk --no-cache add tzdata

# set working directory
WORKDIR /opt/trackdocs

# copy required files from builder
COPY --from=builder /out/trackdocs-api ./trackdocs-api

ENTRYPOINT ["./trackdocs-api"]
