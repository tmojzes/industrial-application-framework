// Copyright 2022 Nokia
// Licensed under the BSD 3-Clause License.
// SPDX-License-Identifier: BSD-3-Clause

package k8sclient

import (
	log "github.com/sirupsen/logrus"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	cfg "sigs.k8s.io/controller-runtime/pkg/client/config"
)

var (
	Config    *rest.Config
	k8sClient client.Client
	scheme    = runtime.NewScheme()
)

func init() {
	utilruntime.Must(corev1.AddToScheme(scheme))      //Scheme for Core V1
}

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
		log.Info("k8sClient already created")
		return k8sClient, nil
	}

	log.Info("no k8sClient")

	if Config == nil {
		log.Info("no config")
		Config = cfg.GetConfigOrDie()
	}

	log.Info("have config")

	cl, err := client.New(Config, client.Options{Scheme: scheme})
	if err != nil {
		log.Error(err, "Failed to create client")
		return nil, err
	}
	log.Info("client created")
	return cl, nil
}
