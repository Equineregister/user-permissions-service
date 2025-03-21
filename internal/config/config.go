package config

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/appconfigdata"
)

// AppConfigParams holds the necessary parameters for AWS AppConfig
type AppConfigParams struct {
	Application   string
	Environment   string
	ConfigProfile string
	PollInterval  time.Duration
}

type Config struct {
	AWSConfig       *aws.Config
	AppConfigClient *appconfigdata.Client
	AppConfigParams AppConfigParams

	Data ConfigData
}

type ConfigData struct {
	LogLevel string `json:"log_level" yaml:"log_level"`
	RDS      RDS    `json:"rds" yaml:"rds"`
}

type RDS struct {
	Database string `json:"database" yaml:"database"`
	Endpoint string `json:"endpoint" yaml:"endpoint"`
	Port     string `json:"port" yaml:"port"`
	User     string `json:"user" yaml:"user"`
}

func setAWS(ctx context.Context) (*aws.Config, error) {
	loadOpts := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithRegion(os.Getenv("AWS_REGION")),
	}

	awsSSOProfile := os.Getenv("AWS_PROFILE")
	if awsSSOProfile != "" {
		slog.Info("Using AWS SSO profile", "profile", awsSSOProfile, "in region", awsSSOProfile)

		loadOpts = append(loadOpts, awsconfig.WithSharedConfigProfile(awsSSOProfile))
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, loadOpts...)
	if err != nil {
		slog.Error("failed to create aws config", "error", err)
		return nil, err
	}

	return &awsCfg, nil
}

// initAppConfig initializes the AWS AppConfig client and parameters
func initAppConfig(_ context.Context, awsCfg *aws.Config) (AppConfigParams, *appconfigdata.Client, error) {
	// Initialize AppConfig client
	appConfigClient := appconfigdata.NewFromConfig(*awsCfg)

	// Set AppConfig parameters from environment variables or use default values
	params := AppConfigParams{
		Application:   getEnvOrDefault("APP_CONFIG_APPLICATION_ID", "user-permissions-service"),
		Environment:   getEnvOrDefault("APP_CONFIG_ENVIRONMENT_ID", ""),
		ConfigProfile: getEnvOrDefault("APP_CONFIG_CONFIGURATION_ID", "default"),
		PollInterval:  time.Duration(getEnvAsIntOrDefault("APP_CONFIG_POLL_INTERVAL_MS", 60000)) * time.Millisecond,
	}

	return params, appConfigClient, nil
}

// getEnvOrDefault returns the value of the environment variable or a default value if not set
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvAsIntOrDefault returns the value of the environment variable as an int or a default value if not set
func getEnvAsIntOrDefault(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	var result int
	_, err := fmt.Sscanf(value, "%d", &result)
	if err != nil {
		slog.Warn("Failed to parse environment variable as int, using default", "key", key, "value", value, "default", defaultValue)
		return defaultValue
	}
	return result
}

// FetchAppConfig fetches the configuration from AWS AppConfig
func (c *Config) FetchAppConfig(ctx context.Context) error {
	// Start a session
	startSessionInput := &appconfigdata.StartConfigurationSessionInput{
		ApplicationIdentifier:                aws.String(c.AppConfigParams.Application),
		ConfigurationProfileIdentifier:       aws.String(c.AppConfigParams.ConfigProfile),
		EnvironmentIdentifier:                aws.String(c.AppConfigParams.Environment),
		RequiredMinimumPollIntervalInSeconds: aws.Int32(int32(c.AppConfigParams.PollInterval.Seconds())),
	}

	session, err := c.AppConfigClient.StartConfigurationSession(ctx, startSessionInput)
	if err != nil {
		return fmt.Errorf("failed to start AppConfig session: %w", err)
	}

	// Get configuration
	getConfigInput := &appconfigdata.GetLatestConfigurationInput{
		ConfigurationToken: session.InitialConfigurationToken,
	}

	resp, err := c.AppConfigClient.GetLatestConfiguration(ctx, getConfigInput)
	if err != nil {
		return fmt.Errorf("failed to get AppConfig configuration: %w", err)
	}

	if len(resp.Configuration) > 0 {
		var configData ConfigData
		if err := json.Unmarshal(resp.Configuration, &configData); err != nil {
			return fmt.Errorf("failed to unmarshal AppConfig configuration: %w", err)
		}
		c.Data = configData
		slog.Info("Successfully fetched configuration from AWS AppConfig")
	} else {
		slog.Warn("Empty configuration received from AWS AppConfig")
	}

	return nil
}

func Load(ctx context.Context) (*Config, error) {
	awsConfig, err := setAWS(ctx)
	if err != nil {
		return nil, fmt.Errorf("set aws: %w", err)
	}

	// Initialize AppConfig
	appConfigParams, appConfigClient, err := initAppConfig(ctx, awsConfig)
	if err != nil {
		return nil, fmt.Errorf("init appconfig: %w", err)
	}

	cfg := &Config{
		AWSConfig:       awsConfig,
		AppConfigClient: appConfigClient,
		AppConfigParams: appConfigParams,
	}

	// Fetch configuration from AppConfig
	if err := cfg.FetchAppConfig(ctx); err != nil {
		slog.Warn("Failed to fetch configuration from AWS AppConfig, continuing with defaults", "error", err)

		return nil, fmt.Errorf("fetch appconfig: %w", err)
	}

	return cfg, nil
}
