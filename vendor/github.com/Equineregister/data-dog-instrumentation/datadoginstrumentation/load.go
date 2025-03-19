package datadoginstrumentation

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

func Load(v *viper.Viper) error {
	if globalConfig != nil {
		return errors.New("Config already set")
	}

	if err := v.BindEnv("datadog.enabled", "DD_ENABLED"); err != nil {
		return err
	}
	if err := v.BindEnv("datadog.profiler.enabled", "DD_PROFILING_ENABLED"); err != nil {
		return err
	}
	if err := v.BindEnv("datadog.profiler.advanced", "DD_PROFILING_ADVANCED"); err != nil {
		return err
	}
	if err := v.BindEnv("datadog.tracer.db.enabled", "DD_ENABLE_DB_TRACER"); err != nil {
		return err
	}
	if err := v.BindEnv("datadog.agentaddr", "DD_AGENT_ADDR"); err != nil {
		return err
	}
	if err := v.BindEnv("service.environment", "SERVICE_ENVIRONMENT"); err != nil {
		return err
	}
	if err := v.BindEnv("service.name", "SERVICE_NAME"); err != nil {
		return err
	}

	// TODO: Add validation rules i.e. can't have profiler disabled but advanced enabled
	globalConfig = &Config{
		Enabled:   v.GetBool("datadog.enabled"),
		Service:   v.GetString("service.name"),
		Env:       v.GetString("service.environment"),
		AgentAddr: v.GetString("datadog.agentaddr"),
		Profiler: ProfilerConfig{
			Enabled:  v.GetBool("datadog.profiler.enabled"),
			Advanced: v.GetBool("datadog.profiler.advanced"),
		},
		Tracer: TracerConfig{
			Enabled: v.GetBool("datadog.enabled"),
		},
		DbTracer: DbTracerConfig{
			Enabled: v.GetBool("datadog.tracer.db.enabled"),
		},
	}
	return nil
}

func getString(varName string) string {
	return os.Getenv(varName)
}

func getBool(varName string) bool {
	return os.Getenv(varName) == "1" || os.Getenv(varName) == "true"
}

func LoadFromEnvVars() error {
	globalConfig = &Config{
		Enabled:   getBool("DD_ENABLED"),
		Service:   getString("SERVICE_NAME"),
		Env:       getString("SERVICE_ENVIRONMENT"),
		AgentAddr: getString("DD_AGENT_ADDR"),
		Profiler: ProfilerConfig{
			Enabled:  getBool("DD_PROFILING_ENABLED"),
			Advanced: getBool("DD_PROFILING_ADVANCED"),
		},
		Tracer: TracerConfig{
			Enabled: getBool("DD_ENABLED"),
		},
		DbTracer: DbTracerConfig{
			Enabled: getBool("DD_ENABLE_DB_TRACER"),
		},
	}

	return nil
}
