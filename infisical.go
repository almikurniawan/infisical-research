package config

import (
	"context"
	"fmt"
	"os"

	infisical "github.com/infisical/go-sdk"
)

// InfisicalClient adalah wrapper untuk Infisical SDK
type InfisicalClient struct {
	client      infisical.InfisicalClientInterface
	projectID   string
	environment string
	secretPath  string
}

// InfisicalConfig adalah konfigurasi untuk membuat InfisicalClient
type InfisicalConfig struct {
	SiteURL      string
	ClientID     string
	ClientSecret string
	ProjectID    string
	Environment  string
	SecretPath   string
}

func NewInfisicalClient(cfg InfisicalConfig) (*InfisicalClient, error) {
	if cfg.SiteURL == "" {
		cfg.SiteURL = "https://app.infisical.com"
	}
	if cfg.SecretPath == "" {
		cfg.SecretPath = "/"
	}
	if cfg.Environment == "" {
		cfg.Environment = "dev"
	}

	client := infisical.NewInfisicalClient(context.Background(), infisical.Config{
		SiteUrl:          cfg.SiteURL,
		AutoTokenRefresh: true,
	})

	_, err := client.Auth().UniversalAuthLogin(cfg.ClientID, cfg.ClientSecret)
	if err != nil {
		return nil, fmt.Errorf("gagal login ke Infisical: %w", err)
	}

	return &InfisicalClient{
		client:      client,
		projectID:   cfg.ProjectID,
		environment: cfg.Environment,
		secretPath:  cfg.SecretPath,
	}, nil
}

func (i *InfisicalClient) GetAllSecrets() (map[string]string, error) {
	secrets, err := i.client.Secrets().List(infisical.ListSecretsOptions{
		ProjectID:          i.projectID,
		Environment:        i.environment,
		SecretPath:         i.secretPath,
		AttachToProcessEnv: false,
	})
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil semua secret: %w", err)
	}

	result := make(map[string]string, len(secrets))
	for _, s := range secrets {
		result[s.SecretKey] = s.SecretValue
	}

	return result, nil
}

func LoadSecret() error {
    // Init Infisical client pakai Machine Identity
    infisical, err := NewInfisicalClient(InfisicalConfig{
        SiteURL:      os.Getenv("INFISICAL_URL"),
        ClientID:     os.Getenv("INFISICAL_CLIENT_ID"),
        ClientSecret: os.Getenv("INFISICAL_CLIENT_SECRET"),
        ProjectID:    os.Getenv("INFISICAL_PROJECT_ID"),
        Environment:  os.Getenv("ENV"), // "prod"
    })
    if err != nil {
        return err
    }

    // Ambil semua secret sekaligus
    secrets, err := infisical.GetAllSecrets()
    if err != nil {
        return err
    }

    // Set ke environment variable
    for key, value := range secrets {
        os.Setenv(key, value)
		fmt.Printf("secret %s %s\n", key, value)
    }

    return nil

}
