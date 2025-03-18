package config

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/Equineregister/data-dog-instrumentation/datadoginstrumentation"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/spf13/viper"
)

type Config struct {
	viper *viper.Viper

	secretsManager *secretsmanager.Client

	AwsConfig *AWSConfig
}

type AWSConfig struct {
	AWS        aws.Config
	SSOProfile string
}

func New() *Config {
	return &Config{
		viper: viper.New(),
	}
}

func (c *Config) GetAWSConfig() aws.Config {
	return c.AwsConfig.AWS
}

func (c *Config) setAWS(ctx context.Context) error {
	if c.AwsConfig != nil {
		// We have already set the AWS Config, so do nothing.
		return nil
	}
	var awscfg aws.Config
	var err error

	ssoProfile := c.viper.GetString("aws.sso.profile")
	region := c.viper.GetString("aws.region")
	if ssoProfile != "" {
		slog.Debug("using aws sso", "profile", ssoProfile, "region", region)

		awscfg, err = awsconfig.LoadDefaultConfig(ctx,
			awsconfig.WithSharedConfigProfile(ssoProfile),
			awsconfig.WithRegion(region),
		)
	} else {
		awscfg, err = awsconfig.LoadDefaultConfig(ctx,
			awsconfig.WithRegion(region),
		)
	}
	if err != nil {
		slog.Error("failed to load AWS config", "error", err)

		return err
	}

	c.AwsConfig = &AWSConfig{
		AWS:        awscfg,
		SSOProfile: c.viper.GetString("aws.sso.profile"),
	}

	c.secretsManager = secretsmanager.NewFromConfig(awscfg)
	return nil
}

type envKey struct {
	key      string
	required bool
}

const (
	LambdaAPI                = "lambda_api"
	LambdaGetUserPermissions = "lambda_get_user_permissions"
)

func (c *Config) Load(ctx context.Context, application string) error {
	appToVarKeys := map[string]map[string]envKey{
		LambdaAPI: {
			"AWS_REGION":             {"aws.region", true},
			"DEFAULT_LANGUAGE_USERS": {"default.language.users", true},
			"AWS_SSO_PROFILE":        {"aws.sso.profile", false},
			"DATABASE_SECRET":        {"db.secret", false},
		},
		LambdaGetUserPermissions: {
			"AWS_REGION":      {"aws.region", true},
			"DATABASE_SECRET": {"db.secret", false},
		},
	}

	varKeys, ok := appToVarKeys[application]
	if !ok {
		return fmt.Errorf("application %s not found", application)
	}

	for envVar, viperKey := range varKeys {
		if viperKey.required && strings.TrimSpace(os.Getenv(envVar)) == "" {
			return fmt.Errorf("%s: required environment variable not set", envVar)
		}

		if err := c.viper.BindEnv(viperKey.key, envVar); err != nil {
			return fmt.Errorf("%s: %w", viperKey.key, err)
		}
	}

	if err := c.setAWS(ctx); err != nil {
		return fmt.Errorf("set aws: %w", err)
	}

	if err := datadoginstrumentation.Load(c.viper); err != nil {
		return fmt.Errorf("load datadog config: %w", err)
	}

	return nil
}
