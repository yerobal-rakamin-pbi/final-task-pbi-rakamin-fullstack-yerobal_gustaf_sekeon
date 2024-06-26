package configbuilder

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	secretManager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"google.golang.org/api/option"
)

type Interface interface {
	BuildConfig()
}

type GCPConfig struct {
	Type        string `json:"type"`
	ProjectID   string `json:"project_id"`
	SecretID    string `json:"-"`
	PrivateKey  string `json:"private_key"`
	ClientEmail string `json:"client_email"`
	SecretName  string `json:"-"`
}

type configbuilder struct {
	gcpConfig      GCPConfig
	configFilePath string
	configJSON     string
}

func Init(gcpConfig GCPConfig, configFilePath string) Interface {
	formattedPrivateKey := regexp.MustCompile(`\\n`).ReplaceAllString(gcpConfig.PrivateKey, "\n")
	gcpConfig.PrivateKey = formattedPrivateKey
	gcpConfig.Type = "service_account"

	return &configbuilder{
		gcpConfig:      gcpConfig,
		configFilePath: configFilePath,
	}
}

func (c *configbuilder) BuildConfig() {
	serviceAccountJSON, err := json.Marshal(c.gcpConfig)
	if err != nil {
		panic(err)
	}

	secretName := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", c.gcpConfig.SecretID, c.gcpConfig.SecretName)

	ctx := context.Background()

	secretManagerClient, err := secretManager.NewClient(ctx, option.WithCredentialsJSON(serviceAccountJSON))
	if err != nil {
		panic(err)
	}

	defer secretManagerClient.Close()

	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretName,
	}

	accessResponse, err := secretManagerClient.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		panic(err)
	}

	c.configJSON = c.beautifyJSON(accessResponse.Payload.Data)
	c.writeConfig()
}

func (c *configbuilder) beautifyJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, data, "", "  ")
	if err != nil {
		panic(err)
	}

	return prettyJSON.String()
}

func (c *configbuilder) writeConfig() {
	// Write the config to the file, if the file doesn't exist, create it. If it exists, overwrite it.
	configFile, err := os.Create(c.configFilePath)
	if err != nil {
		panic(err)
	}

	defer configFile.Close()

	_, err = configFile.WriteString(c.configJSON)
	if err != nil {
		panic(err)
	}
}
