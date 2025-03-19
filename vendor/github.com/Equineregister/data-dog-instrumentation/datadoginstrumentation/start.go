package datadoginstrumentation

import (
	"log/slog"
	"net"
	"net/http"
	"time"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

var stops = make([]func(), 0)

/*
The configuration used here is the configuration defined as the default handler
used by the datadog package. The one change is upping the time timeout from 2s->5s
*/

var defaultDialer = &net.Dialer{
	Timeout:   30 * time.Second,
	KeepAlive: 30 * time.Second,
	DualStack: true,
}

var defaultClient = &http.Client{
	// We copy the transport to avoid using the default one, as it might be
	// augmented with tracing and we don't want these calls to be recorded.
	// See https://golang.org/pkg/net/http/#DefaultTransport .
	Transport: &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           defaultDialer.DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
	Timeout: 5 * time.Second,
}

type ddLogger struct {
	log *slog.Logger
}

func (l ddLogger) Log(msg string) {
	l.log.Warn(msg)
}

func Start() error {
	if globalConfig == nil || !globalConfig.Enabled {
		return nil
	}

	if globalConfig.Tracer.Enabled {
		logger := ddLogger{log: slog.Default()}
		tracer.Start(
			tracer.WithAgentAddr(globalConfig.AgentAddr),
			tracer.WithEnv(globalConfig.Env),
			tracer.WithService(globalConfig.Service),
			tracer.WithHTTPClient(defaultClient),
			tracer.WithLogger(logger),
		)
		stops = append(stops, tracer.Stop)
	}

	if globalConfig.Profiler.Enabled {
		profilers := []profiler.ProfileType{
			profiler.CPUProfile,
			profiler.HeapProfile,
		}
		if globalConfig.Profiler.Advanced {
			profilers = append(profilers, profiler.BlockProfile)
			profilers = append(profilers, profiler.MutexProfile)
			profilers = append(profilers, profiler.GoroutineProfile)
		}
		err := profiler.Start(
			profiler.WithEnv(globalConfig.Env),
			profiler.WithAgentAddr(globalConfig.AgentAddr),
			profiler.WithService(globalConfig.Service),
			profiler.WithProfileTypes(profilers...),
		)
		if err != nil {
			return err
		}
		stops = append(stops, profiler.Stop)
	}

	return nil
}

func Stop() {
	for _, stop := range stops {
		stop()
	}
}
