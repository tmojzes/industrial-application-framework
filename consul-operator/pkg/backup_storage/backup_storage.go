package backupStorage

import (
	"context"
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/nokia/industrial-application-framework/consul-operator/pkg/template"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/json"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"strings"
)

var log = logf.Log.WithName("controller_consul")

type BackupCRStatus struct {
	client.Client
}

var s3Endpoint, bucketName, accessKey, secretAccessKey string

var BackupCRStat *BackupCRStatus

func (bStatus *BackupCRStatus) createMinioClient() (*minio.Client, error) {
	logger := log.WithName("backup_storage").WithName("createMinioClient")
	logger.Info("Called")
	minioClient, err := minio.New(s3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to initialize Minio Client, accessKey: %v", accessKey)
	}
	log.Info("Minio Client is initialized", "accessKey", accessKey)
	return minioClient, nil
}

func (bStatus *BackupCRStatus) uploadFile(cl *minio.Client, content string) error {
	logger := log.WithName("backup_storage").WithName("uploadFile")
	logger.Info("Called")

	filePath := strings.Join([]string{os.Getenv(template.DeploymentDir), "consul_content.txt"}, "/")

	file, err := os.Create(filePath)
	if err != nil {
		log.Error(err, "creating file")
		return err
	}
	defer file.Close()
	//	defer os.Remove(file.Name())

	if _, err = file.WriteString(content); err != nil {
		log.Error(err, "writing into file")
		return err
	}

	err = file.Sync()
	if err != nil {
		log.Error(err, "syncing file")
		return err
	}

	fileStat, err := file.Stat()
	if err != nil {
		log.Error(err, "creating file stat")
		return err
	}

	log.Info("file", "size", fileStat.Size())

	uploadInfo, err := cl.FPutObject(context.Background(), bucketName, file.Name(), filePath, minio.PutObjectOptions{ContentType: ""})
	if err != nil {
		log.Error(err, "putting file into bucket")
		return err
	}
	fmt.Println("--------Client Successfully uploaded object: ", uploadInfo)

	return nil
}

func (bStatus *BackupCRStatus) getBackupCRStatus(nameSpace string) bool {

	logger := log.WithName("backupCRStatus").WithName("getBackupCRStatus")
	logger.Info("Called")
	BackupCRs := &unstructured.UnstructuredList{}
	BackupCRs.SetGroupVersionKind(schema.GroupVersionKind{Group: "ops.dac.nokia.com", Version: "v1alpha1", Kind: "Backup"})
	err := bStatus.Client.List(context.TODO(), BackupCRs, &client.ListOptions{Namespace: nameSpace})

	if err != nil {
		logger.Error(err, "Error listing Backup CRs")
		return false
	}

	for _, item := range BackupCRs.Items {
		field, found, err := unstructured.NestedMap(item.Object, "metadata")
		if err == nil && found && field["name"].(string) == "backup-consul" {
			field, found, err := unstructured.NestedMap(item.Object, "status")
			if err == nil && found {
				s3Endpoint = field["s3Endpoint"].(string)
				bucketName = field["bucketConfiguration"].(map[string]interface{})["bucketName"].(string)
				accessKey = field["bucketConfiguration"].(map[string]interface{})["accessKey"].(string)
				secretAccessKey = field["bucketConfiguration"].(map[string]interface{})["secretAccessKey"].(string)
				logger.Info("backup CR found")
				log.Info("bucket", "s3Endpoint", s3Endpoint)
				log.Info("bucket", "bucketName", bucketName)
				log.Info("bucket", "accessKey", accessKey)
				log.Info("bucket", "secretAccessKey", secretAccessKey)

				return true
			}
			return false
		}
	}
	logger.Info("exit")
	return false
}

func (bStatus *BackupCRStatus) readConsulContent() (string, error) {
	logger := log.WithName("backup_storage").WithName("readConsulContent")
	logger.Info("Called")

	conf := consulapi.DefaultConfig()
	conf.Address = fmt.Sprintf("consul.default.svc.cluster.local:8500")
	consulClient, err := consulapi.NewClient(conf)
	if err != nil {
		logger.Error(err, "creating consul api client")
		return "", err
	}

	KVPairs, _, err := consulClient.KV().List("/", nil)
	if err != nil {
		logger.Error(err, "listing consul content")
		return "", err
	}
	log.Info("consul content", "KVPairs", KVPairs)

	consulContent, err := json.Marshal(KVPairs)
	if err != nil {
		logger.Error(err, "marshal the KVPairs map")
		return "", err
	}

	return string(consulContent), nil
}

func (bStatus *BackupCRStatus) UploadDataToBucket(nameSpace string) {
	logger := log.WithName("backup_storage").WithName("UploadDataToBucket")
	logger.Info("Called")

	if bStatus.getBackupCRStatus(nameSpace) {
		minioClient, err := bStatus.createMinioClient()
		if err != nil {
			return
		}
		consulContent, err := bStatus.readConsulContent()
		if err != nil {
			return
		}
		err = bStatus.uploadFile(minioClient, consulContent)
		if err != nil {
			return
		}
	}

	return
}
