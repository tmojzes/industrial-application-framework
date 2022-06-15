# Copyright 2022 Nokia
# Licensed under the BSD 3-Clause License.
# SPDX-License-Identifier: BSD-3-Clause

# Build the manager binary
FROM golang:1.16 as builder

WORKDIR /workspace

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

# Copy the go source
COPY main.go main.go
COPY pkg/ pkg/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o consul-backup main.go

FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

ENV DEPLOYMENT_DIR=/usr/src/app

WORKDIR /
COPY --from=builder /workspace/consul-backup .

USER 65532:65532

ENTRYPOINT ["/consul-backup"]
