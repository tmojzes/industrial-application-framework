// Copyright 2022 Nokia
// Licensed under the BSD 3-Clause License.
// SPDX-License-Identifier: BSD-3-Clause

package consulclient

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/nokia/industrial-application-framework/consul-backup/pkg/serviceconfig"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/json"
)

func CreateConsulClient() (*consulapi.Client, error) {
	log.Info("CreateConsulClient called")

	conf := consulapi.DefaultConfig()
	conf.Address = fmt.Sprintf(serviceconfig.ConfigData.ConsulAddress)

	consulClient, err := consulapi.NewClient(conf)
	if err != nil {
		log.Error(err, "Failed to creat consul api client")
		return nil, err
	}
	return consulClient, nil
}

func ReadConsulContent(consulClient *consulapi.Client) (string, error) {
	log.Info("ReadConsulContent called")

	KVPairs, _, err := consulClient.KV().List("/", nil)
	if err != nil {
		log.Error(err, "Failed to list consul content")
		return "", err
	}
	log.Info("consul content", "KVPairs", KVPairs)

	consulContent, err := json.Marshal(KVPairs)
	if err != nil {
		log.Error(err, "Failed to marshal the KVPairs map")
		return "", err
	}

	return string(consulContent), nil
}
