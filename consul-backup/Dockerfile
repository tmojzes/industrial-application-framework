# Copyright 2022 Nokia
# Licensed under the BSD 3-Clause License.
# SPDX-License-Identifier: BSD-3-Clause

# Build the manager binary
FROM golang:1.17 as builder

WORKDIR /workspace

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY pkg/ pkg/
COPY Makefile Makefile

# Build
RUN make bin-build

FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

ENV DEPLOYMENT_DIR=/usr/src/app

WORKDIR /
COPY --from=builder /workspace/consul-backup .

USER 65532:65532

ENTRYPOINT ["/consul-backup"]
