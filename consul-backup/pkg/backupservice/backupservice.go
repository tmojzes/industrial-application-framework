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

func (handler *BackupCRHandler) fillBackupStoreAccessInfo(nameSpace string) error {
	log.Info("getBackupCRHandler called")

	BackupCR := &unstructured.Unstructured{}
	BackupCR.SetGroupVersionKind(schema.GroupVersionKind{Group: "ops.dac.nokia.com", Version: "v1alpha1", Kind: "Backup"})
	err := handler.Client.Get(context.TODO(), client.ObjectKey{ Namespace: nameSpace, Name: serviceconfig.ConfigData.BackupCrName}, BackupCR)

	if err != nil {
		return errors.Wrap(err, "Failed to get backup CRs")
	}

	field, found, err := unstructured.NestedMap(BackupCR.Object, "status")
	if err != nil {
		return errors.Wrap(err, "Failed to read backup CR")
	}
	
	if found {
		handler.s3Endpoint = field["s3Endpoint"].(string)
		handler.bucketName = field["bucketConfiguration"].(map[string]interface{})["bucketName"].(string)
		handler.accessKey = field["bucketConfiguration"].(map[string]interface{})["accessKey"].(string)
		handler.secretAccessKey = field["bucketConfiguration"].(map[string]interface{})["secretAccessKey"].(string)
		log.Info("backup CR found")
		return nil
	} else {
		return errors.New("Status not found in backup CR")
	}

}

func (handler *BackupCRHandler) uploadDataToBackupStore (nameSpace string) error {
	log.Info("UploadDataToBucket called")

	err := handler.fillBackupStoreAccessInfo(nameSpace)
	if err != nil {
		return errors.Wrap(err, "Failed to get backup store access info")
	}
	consulClient, err := consulclient.CreateConsulClient()
	if err != nil {
		return errors.Wrap(err, "Fail to create consul client")
	}
	consulContent, err := consulclient.ReadConsulContent(consulClient)
	if err != nil {
		return errors.Wrap(err, "Failed to read consul content")
	}
	s3Cl, err := s3client.CreateS3Client(handler.s3Endpoint, handler.accessKey, handler.secretAccessKey)
	if err != nil {
		return errors.Wrap(err, "Failed to create backup store client")
	}
	err = s3Cl.UploadFileToS3Storage( consulContent, handler.bucketName)
	if err != nil {
		return errors.Wrap(err, "Failed to store in backup store")
	}

	return nil
}

func StartPeriodicBackup(nameSpace string) error {
	log.Info("BackupService called")

	k8sClient, err := k8sclient.GetK8sClient()
	if err != nil {
		return errors.Wrap(err, "Failed to get k8s client")
	}

	handler := &BackupCRHandler{Client: k8sClient}
	err = handler.uploadDataToBackupStore(nameSpace)
	if err != nil {
		log.Error(err, "Failed to upload data to backup storage")
	}

	duration, err := time.ParseDuration(serviceconfig.ConfigData.Duration)
	if err != nil {
		return errors.Wrap(err, "Failed to parse duration")
	}

	ticker := time.NewTicker(duration)
	for range ticker.C {
		err = handler.uploadDataToBackupStore(nameSpace)
		if err != nil {
			log.Error(err, "Failed to upload a file to backup storage")
		}
		log.Info("sleeping...")
	}

	return errors.New("Periodical backup exited")
}
