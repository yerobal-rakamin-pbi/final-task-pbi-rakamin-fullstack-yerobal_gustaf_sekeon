package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"regexp"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
	"rakamin-final-task/helpers/files"
)

type GCPConfig struct {
	Type        string `json:"type"`
	ProjectID   string `json:"project_id"`
	PrivateKey  string `json:"private_key"`
	ClientEmail string `json:"client_email"`
}

type storageLib struct {
	Config     GCPConfig
	BucketName string
	client     *storage.Client
}

type Interface interface {
	Upload(ctx context.Context, file *files.File, path string) (string, error)
	UploadFromBytes(ctx context.Context, file *bytes.Reader, fileName string, path string) (string, error)
	Delete(ctx context.Context, path string, fileName string) error
	getObjectPlace(objectPath string) *storage.ObjectHandle
}

func Init(config GCPConfig, bucketName string) Interface {
	config.PrivateKey = regexp.MustCompile(`\\n`).ReplaceAllString(config.PrivateKey, "\n")
	config.Type = "service_account"

	serviceAccountJSON, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(serviceAccountJSON))
	if err != nil {
		panic(err)
	}

	return &storageLib{
		Config:     config,
		BucketName: bucketName,
		client:     client,
	}
}

func (s *storageLib) getObjectPlace(objectPath string) *storage.ObjectHandle {
	return s.client.Bucket(s.BucketName).Object(objectPath)
}

func (s *storageLib) Upload(ctx context.Context, file *files.File, path string) (string, error) {
	var imageURL string
	writer := s.getObjectPlace(fmt.Sprintf("%s/%s", path, file.Meta.Filename)).NewWriter(ctx)

	if _, err := io.Copy(writer, file.Content); err != nil {
		return imageURL, err
	}

	if err := writer.Close(); err != nil {
		return imageURL, err
	}

	parsedURL, err := url.Parse(fmt.Sprintf("https://storage.googleapis.com/%s/%s/%s", s.BucketName, path, file.Meta.Filename))
	if err != nil {
		return imageURL, err
	}

	imageURL = parsedURL.String()

	return imageURL, nil
}

func (s *storageLib) UploadFromBytes(ctx context.Context, file *bytes.Reader, fileName string, path string) (string, error) {
	var imageURL string
	writer := s.getObjectPlace(path + "/" + fileName).NewWriter(ctx)

	if _, err := io.Copy(writer, file); err != nil {
		return imageURL, err
	}

	if err := writer.Close(); err != nil {
		return imageURL, err
	}

	parsedURL, err := url.Parse(fmt.Sprintf("https://storage.googleapis.com/%s/%s/%s", s.BucketName, path, fileName))
	if err != nil {
		return imageURL, err
	}

	imageURL = parsedURL.String()

	return imageURL, nil
}

func (s *storageLib) Delete(ctx context.Context, filename string, path string) error {
	return s.getObjectPlace(path + "/" + filename).Delete(ctx)
}
