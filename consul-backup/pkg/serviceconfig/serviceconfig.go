// Copyright 2022 Nokia
// Licensed under the BSD 3-Clause License.
// SPDX-License-Identifier: BSD-3-Clause

package serviceconfig

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Configuration struct {
	Duration        string `yaml:"duration"`
	ConsulAddress   string `yaml:"consulAddress"`
	BackupCrName    string `yaml:"backupCrName"`
}

const configFileKey = "CONFIG_FILE"

var ConfigData Configuration

func ReadServiceConfig() error {
	log.Info("ReadServiceConfig called")

	configFile, err := os.Open(os.Getenv(configFileKey))
	if err != nil {
		log.Error(err, "Failed to open config file")
		return err
	}

	byteValue, err := ioutil.ReadAll(configFile)
	if err != nil {
		log.Error(err, "Failed to read from config file")
		return err
	}

	err = yaml.Unmarshal(byteValue, &ConfigData)

	if err != nil {
		log.Error(err, "Failed to unmarshal the config data")
		return err
	}

	return nil
}
