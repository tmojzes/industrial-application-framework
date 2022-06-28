// Copyright 2022 Nokia
// Licensed under the BSD 3-Clause License.
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"fmt"
	backupService "github.com/nokia/industrial-application-framework/consul-backup/pkg/backupservice"
	"github.com/nokia/industrial-application-framework/consul-backup/pkg/serviceconfig"
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

	err = serviceconfig.ReadServiceConfig()
	if err != nil {
		log.Error(err, "Failed to read configmap ")
		panic(err)
	}

	err = backupService.StartPeriodicBackup(ownNamespace)
	if err != nil {
		log.Error(err, "Failed to start periodical backup ")
		panic(err)
	}
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
