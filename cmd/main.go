package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"gopkg.in/yaml.v3"
)

type Metadata struct {
	Name        string            `yaml:"name"`
	Namespace   string            `yaml:"namespace,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
	Labels      map[string]string `yaml:"labels,omitempty"`
}

type Secret struct {
	APIVersion string            `yaml:"apiVersion"`
	Kind       string            `yaml:"kind"`
	Type       string            `yaml:"type"`
	Metadata   Metadata          `yaml:"metadata"`
	Data       map[string]string `yaml:"data"`
}

type SecretName struct {
	Name string `yaml:"name"`
	Key  string `yaml:"key,omitempty"`
}

type SecretManager struct {
	APIVersion string       `yaml:"apiVersion"`
	Kind       string       `yaml:"kind"`
	Metadata   Metadata     `yaml:"metadata"`
	Secrets    []SecretName `yaml:"secrets"`
}

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatalf("Expected exactly 1 argument, got: %q", flag.Args())
	}
	generator := flag.Arg(0)

	projectId := os.Getenv("PROJECT_ID")
	if projectId == "" {
		log.Fatal("Expected env PROJECT_ID to be set")
	}

	in, err := ioutil.ReadFile(generator)
	if err != nil {
		log.Fatalf("Failed to read generator manifest: %s", err)
	}

	var parsed SecretManager
	if err := yaml.Unmarshal(in, &parsed); err != nil {
		log.Fatalf("Invalid generator manifest: %s", err)
	}

	secretData := map[string]string{}
	for _, secret := range parsed.Secrets {
		name := "projects/" + projectId + "/secrets/" + secret.Name + "/versions/latest"
		value, err := accessSecretVersion(name)
		if err != nil {
			log.Fatalf("Failed to load secret from secret manager: %s", err)
		}
		key := secret.Key
		if key == "" {
			key = secret.Name
		}
		secretData[key] = base64.StdEncoding.EncodeToString(value)
	}

	manifest := Secret{
		APIVersion: "v1",
		Kind:       "Secret",
		Metadata: Metadata{
			Name:        parsed.Metadata.Name,
			Namespace:   parsed.Metadata.Namespace,
			Annotations: parsed.Metadata.Annotations,
			Labels:      parsed.Metadata.Labels,
		}, Type: "Opaque",
		Data: secretData,
	}

	b, err := yaml.Marshal(manifest)
	if err != nil {
		log.Fatalf("Failed to marshal secret: %s", err)
	}

	fmt.Print(string(b))
}

func accessSecretVersion(name string) ([]byte, error) {
	// name := "projects/my-project/secrets/my-secret/versions/5"
	// name := "projects/my-project/secrets/my-secret/versions/latest"

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create secretmanager client: %v", err)
	}

	// Build the request.
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to access secret version: %v", err)
	}

	return result.Payload.Data, nil
}
