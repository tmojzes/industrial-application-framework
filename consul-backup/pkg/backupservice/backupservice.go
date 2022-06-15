package backupService

import (
	"context"
	"github.com/nokia/industrial-application-framework/consul-backup/pkg/consulclient"
	"github.com/nokia/industrial-application-framework/consul-backup/pkg/k8sclient"
	"github.com/nokia/industrial-application-framework/consul-backup/pkg/minioclient"
	"time"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type BackupCRStatus struct {
	client.Client
}

var s3Endpoint, bucketName, accessKey, secretAccessKey string

var BackupCRStat *BackupCRStatus

func (bStatus *BackupCRStatus) getBackupCRStatus(nameSpace string) bool {
	log.Info("getBackupCRStatus called")

	BackupCR := &unstructured.Unstructured{}
	BackupCR.SetGroupVersionKind(schema.GroupVersionKind{Group: "ops.dac.nokia.com", Version: "v1alpha1", Kind: "Backup"})
    err := bStatus.Client.Get(context.TODO(), client.ObjectKey{ Namespace: nameSpace, Name: "backup-consul"}, BackupCR)

	if err != nil {
		log.Error(err, "Failed to get backup CRs")
		return false
	}

	field, found, err := unstructured.NestedMap(BackupCR.Object, "status")
	if err == nil && found {
		s3Endpoint = field["s3Endpoint"].(string)
		bucketName = field["bucketConfiguration"].(map[string]interface{})["bucketName"].(string)
		accessKey = field["bucketConfiguration"].(map[string]interface{})["accessKey"].(string)
		secretAccessKey = field["bucketConfiguration"].(map[string]interface{})["secretAccessKey"].(string)
		log.Info("backup CR found")
		return true
	}

	log.Error("Error reading backup CR status")
	return false
}

func (bStatus *BackupCRStatus) uploadDataToBucket(nameSpace string) {
	log.Info("UploadDataToBucket called")

	if bStatus.getBackupCRStatus(nameSpace) {
		minioClient, err := minioclient.CreateMinioClient(s3Endpoint, accessKey, secretAccessKey)
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
		err = minioclient.UploadFileToMinio(minioClient, consulContent, bucketName)
		if err != nil {
			return
		}
	}

	return
}

func BackupService(nameSpace string) {
	log.Info("BackupService called")

	k8sclient, err := k8sclient.GetK8sClient()

	if err != nil {
		log.Error(err, "Get client error")
	} else {
		BackupCRStat = &BackupCRStatus{Client: k8sclient}
		for {
			BackupCRStat.uploadDataToBucket(nameSpace)
			log.Infof("sleeping...")
			time.Sleep(1*time.Hour)
		}
	}
}
