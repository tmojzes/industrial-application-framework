# Copyright 2022 Nokia
# Licensed under the BSD 3-Clause License.
# SPDX-License-Identifier: BSD-3-Clause

bin-build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -v -o  consul-backup  main.go

docker-build: bin-build
	docker build --tag consul-backup .
