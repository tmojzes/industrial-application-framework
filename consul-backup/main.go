// Copyright 2020 Nokia
// Licensed under the BSD 3-Clause License.
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"fmt"
	backupStorage "github.com/nokia/industrial-application-framework/consul-backup/pkg/backup_example"
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/nokia/industrial-application-framework/consul-backup/pkg/k8sclient"

	"os"
)


func main() {
	log.Infof("Starting consul-backup")
	watchNamespace, err := getWatchNamespace()
	if err != nil {
		log.Error(err, "unable to get WatchNamespace, "+
			"the manager will watch and manage resources in all namespaces")
	}

/*	cl, err := client.New(cfg.GetConfigOrDie(), client.Options{Scheme: scheme})
	if err != nil {
		fmt.Println("failed to create client")
		os.Exit(1)
	}
 */
	k8sclient, err := k8sclient.GetK8sClient()
	log.Infof("after create client")

	if err != nil {
		log.Error(err, "get client error")
	} else {
		backupStorage.BackupCRStat = &backupStorage.BackupCRStatus{Client: k8sclient}
		log.Infof("before upload")
		backupStorage.BackupCRStat.UploadDataToBucket(watchNamespace)
	}

	for {
		log.Infof("sleeping...")
		time.Sleep(1*time.Hour)
	}

}

// getWatchNamespace returns the Namespace the operator should be watching for changes
func getWatchNamespace() (string, error) {
	// WatchNamespaceEnvVar is the constant for env variable WATCH_NAMESPACE
	// which specifies the Namespace to watch.
	// An empty value means the operator is running with cluster scope.
	log.Infof("get namespace")
	var watchNamespaceEnvVar = "WATCH_NAMESPACE"

	ns, found := os.LookupEnv(watchNamespaceEnvVar)
	if !found {
		return "", fmt.Errorf("%s must be set", watchNamespaceEnvVar)
	}
	return ns, nil
}
