// Copyright 2022 Nokia
// Licensed under the BSD 3-Clause License.
// SPDX-License-Identifier: BSD-3-Clause

package s3client

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

type Client struct {
	*minio.Client
}

const workDir = "/tmp"

func CreateS3Client(s3Endpoint, accessKey, secretAccessKey string) (*Client, error) {
	log.Info("createMinioClient called")
	minioClient, err := minio.New(s3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to initialize Minio Client, accessKey: %v", accessKey)
	}
	log.Info("minio Client is initialized", "accessKey", accessKey)
	return &Client{minioClient}, nil
}

func (cl Client) UploadFileToS3Storage(content, bucketName string) error {
	log.Info("uploadFile called")

	filePath := strings.Join([]string{workDir, "consulContent.json"}, "/")

	file, err := os.Create(filePath)
	if err != nil {
		return errors.Wrap(err, "Failed to create file")
	}
	defer file.Close()

	if _, err = file.WriteString(content); err != nil {
		return errors.Wrap(err, "Failed to write file")
	}

	err = file.Sync()
	if err != nil {
		return errors.Wrap(err, "Failed to sync file")
	}

	// In this function a file is uploaded with the same name to the bucket, it overwrites the previous version of the file.
	// If the earlier version of the file is needed different file name shall be used.
	uploadInfo, err := cl.FPutObject(context.Background(), bucketName, file.Name(), filePath, minio.PutObjectOptions{})
	if err != nil {
		return errors.Wrap(err, "Failed to put file into bucket")
	}
	log.Info("--------client successfully", "uploaded object:", uploadInfo)

	return nil
}
