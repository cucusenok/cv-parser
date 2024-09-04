# syntax=docker/dockerfile:1

FROM golang:1.20.6 as builder

RUN apt update
WORKDIR app


RUN --mount=type=cache,target=/go/pkg/mod/ \
      --mount=type=bind,source=./app/go.sum,target=go.sum \
      --mount=type=bind,source=./app/go.mod,target=go.mod \
      go mod download -x
