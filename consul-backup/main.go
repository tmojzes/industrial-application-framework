// Copyright 2022 Nokia
// Licensed under the BSD 3-Clause License.
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"fmt"
	backupService "github.com/nokia/industrial-application-framework/consul-backup/pkg/backupservice"
	log "github.com/sirupsen/logrus"

	"os"
)


func main() {
	log.Infof("Starting consul-backup")
	ownNamespace, err := getOwnNamespace()
	if err != nil {
		log.Error(err, "Unable to get OwnNamespace ")
		panic(err)
	}

	backupService.BackupService(ownNamespace)
}

func getOwnNamespace() (string, error) {
	log.Infof("get namespace")
	var ownNamespace = "OWN_NAMESPACE"

	ns, found := os.LookupEnv(ownNamespace)
	if !found {
		return "", fmt.Errorf("%s must be set", ownNamespace)
	}
	return ns, nil
}
