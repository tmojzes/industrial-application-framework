# Copyright 2022 Nokia
# Licensed under the BSD 3-Clause License.
# SPDX-License-Identifier: BSD-3-Clause

PACKAGE=consul-backup

bin-build:
	CGO_ENABLED=0 GOOS=linux go build -a -v -o  ${PACKAGE}  main.go

docker-build: bin-build
	docker build --tag consul-backup .
