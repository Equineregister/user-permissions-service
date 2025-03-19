package datadoginstrumentation

import (
	"net/http"

	lambdatrace "github.com/DataDog/datadog-lambda-go"
	"google.golang.org/grpc"
	grpctrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
)

var namer = func(r *http.Request) string {
	return r.URL.Path
}

func WrapHttp(handler http.Handler) http.Handler {
	if globalConfig == nil || !globalConfig.Enabled {
		return handler
	}
	if globalConfig.Tracer.Enabled {
		handler = httptrace.WrapHandler(
			handler,
			globalConfig.Service,
			"http",
			httptrace.WithResourceNamer(namer))
	}
	return handler
}

func WrapGrpc(mdlwrs []grpc.UnaryServerInterceptor) []grpc.UnaryServerInterceptor {
	if globalConfig == nil || !globalConfig.Enabled {
		return mdlwrs
	}

	if globalConfig.Tracer.Enabled {
		grpcmdlwr := grpctrace.UnaryServerInterceptor(
			grpctrace.WithServiceName(globalConfig.Service),
		)
		mdlwrs = append(mdlwrs, grpcmdlwr)
	}

	return mdlwrs
}

// This is possible but requires a big change in how database connections are
// initialised. This is a bigger change than is required for now so going to
// leave this function here and this link: https://pkg.go.dev/gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql
func WrapDB() {}

func GrpcClientInterceptor() grpc.UnaryClientInterceptor {
	return grpctrace.UnaryClientInterceptor(
		grpctrace.WithServiceName(globalConfig.Service),
	)
}

func WrapGrpcClient(mddlwr []grpc.UnaryClientInterceptor) []grpc.UnaryClientInterceptor {
	if globalConfig == nil || !globalConfig.Enabled {
		return mddlwr
	}
	return append(mddlwr, GrpcClientInterceptor())
}

var config = &lambdatrace.Config{
	DDTraceEnabled:  true,
	EnhancedMetrics: true,
}

func WrapLambdaHandler(handler interface{}) interface{} {
	return lambdatrace.WrapFunction(handler, config)
}
