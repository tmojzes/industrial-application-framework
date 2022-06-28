// Copyright 2022 Nokia
// Licensed under the BSD 3-Clause License.
// SPDX-License-Identifier: BSD-3-Clause

package k8sclient

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	cfg "sigs.k8s.io/controller-runtime/pkg/client/config"
)

var (
	Config    *rest.Config
	k8sClient client.Client
	scheme    = runtime.NewScheme()
)

func GetK8sClient() (client.Client, error) {
	if k8sClient != nil {
		return k8sClient, nil
	} else {
		return CreateK8sClient()
	}
}

func CreateK8sClient() (client.Client, error) {
	log.Info("createK8sClient called")

	if k8sClient != nil {
		return k8sClient, nil
	}

	if Config == nil {
		Config = cfg.GetConfigOrDie()
	}

	cl, err := client.New(Config, client.Options{Scheme: scheme})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create client")
	}

	return cl, nil
}
