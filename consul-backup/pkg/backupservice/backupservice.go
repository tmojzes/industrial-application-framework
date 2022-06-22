// Copyright 2022 Nokia
// Licensed under the BSD 3-Clause License.
// SPDX-License-Identifier: BSD-3-Clause

package backupService

import (
	"context"
	"github.com/nokia/industrial-application-framework/consul-backup/pkg/consulclient"
	"github.com/nokia/industrial-application-framework/consul-backup/pkg/k8sclient"
	"github.com/nokia/industrial-application-framework/consul-backup/pkg/s3client"
	"github.com/nokia/industrial-application-framework/consul-backup/pkg/serviceconfig"
	"github.com/pkg/errors"
	"time"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type BackupCRHandler struct {
	client.Client
	s3Endpoint string
	bucketName string
	accessKey string
	secretAccessKey string
}

var handler *BackupCRHandler

func (handler *BackupCRHandler) fillS3AccessInfo(nameSpace string) error {
	log.Info("getBackupCRHandler called")

	BackupCR := &unstructured.Unstructured{}
	BackupCR.SetGroupVersionKind(schema.GroupVersionKind{Group: "ops.dac.nokia.com", Version: "v1alpha1", Kind: "Backup"})
	err := handler.Client.Get(context.TODO(), client.ObjectKey{ Namespace: nameSpace, Name: serviceconfig.ConfigData.BackupCrName}, BackupCR)

	if err != nil {
		log.Error(err, "Failed to get backup CRs")
		return err
	}

	field, found, err := unstructured.NestedMap(BackupCR.Object, "status")
	if err != nil {
		log.Error(err, "Failed to read backup CR")
		return err
	}
	
	if found {
		handler.s3Endpoint = field["s3Endpoint"].(string)
		handler.bucketName = field["bucketConfiguration"].(map[string]interface{})["bucketName"].(string)
		handler.accessKey = field["bucketConfiguration"].(map[string]interface{})["accessKey"].(string)
		handler.secretAccessKey = field["bucketConfiguration"].(map[string]interface{})["secretAccessKey"].(string)
		log.Info("backup CR found")
		return nil
	} else {
		log.Error("Status not found in backup CR")
		return errors.New("Status not found in backup CR")
	}

}

func (handler *BackupCRHandler) uploadDataToS3Storage (nameSpace string) {
	log.Info("UploadDataToBucket called")

	err := handler.fillS3AccessInfo(nameSpace)
	if err != nil {
		return
	}
	consulClient, err := consulclient.CreateConsulClient()
	if err != nil {
		return
	}
	consulContent, err := consulclient.ReadConsulContent(consulClient)
	if err != nil {
		return
	}
	s3Cl, err := s3client.CreateS3Client(handler.s3Endpoint, handler.accessKey, handler.secretAccessKey)
	if err != nil {
		return
	}
	err = s3Cl.UploadFileToS3Storage( consulContent, handler.bucketName)
	if err != nil {
		return
	}

	return
}

func StartPeriodicBackup(nameSpace string) {
	log.Info("BackupService called")

	err := serviceconfig.ReadServiceConfig()
	if err != nil {
		return
	}

	k8sClient, err := k8sclient.GetK8sClient()
	if err != nil {
		return
	}

	handler = &BackupCRHandler{Client: k8sClient}
	handler.uploadDataToS3Storage(nameSpace)

	duration, err := time.ParseDuration(serviceconfig.ConfigData.Duration)
	if err != nil {
		log.Error(err, "Failed to parse duration")
		return
	}

	ticker := time.NewTicker(duration)
	for range ticker.C {
		err = serviceconfig.ReadServiceConfig()
		if err == nil {
			handler.uploadDataToS3Storage(nameSpace)
			log.Info("sleeping...")
		}
	}
}
