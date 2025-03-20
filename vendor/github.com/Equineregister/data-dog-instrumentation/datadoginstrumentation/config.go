package datadoginstrumentation

var globalConfig *Config

type Config struct {
	Enabled   bool
	Service   string
	Env       string
	AgentAddr string
	Profiler  ProfilerConfig
	Tracer    TracerConfig
	DbTracer  DbTracerConfig
}

type ProfilerConfig struct {
	Enabled  bool
	Advanced bool
}

type TracerConfig struct {
	Enabled bool
}

type DbTracerConfig struct {
	Enabled bool
}

func DbTracerEnabled() bool {
	if globalConfig == nil {
		return false
	}

	return globalConfig.DbTracer.Enabled
}

func ServiceName() string {
	if globalConfig == nil {
		return ""
	}
	return globalConfig.Service
}

func Enabled() bool {
	return globalConfig != nil && globalConfig.Enabled
}
