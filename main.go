// Copyright (c) 2024 AccelByte Inc. All Rights Reserved.
// This is licensed software from AccelByte Inc, for limitations
// and restrictions contact your company contract manager.

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"accelbyte.net/session-manager-grpc-plugin-server-go/pkg/common"
	"accelbyte.net/session-manager-grpc-plugin-server-go/pkg/server"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	registered_v1 "accelbyte.net/session-manager-grpc-plugin-server-go/pkg/pb"
	"github.com/AccelByte/accelbyte-go-sdk/services-api/pkg/factory"
	"github.com/AccelByte/accelbyte-go-sdk/services-api/pkg/repository"
	"github.com/AccelByte/accelbyte-go-sdk/services-api/pkg/service/iam"
	sdkAuth "github.com/AccelByte/accelbyte-go-sdk/services-api/pkg/utils/auth"
	promgrpc "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/prometheus/client_golang/prometheus"
	prometheusCollectors "github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

const (
	environment     = "production"
	id              = int64(1)
	metricsEndpoint = "/metrics"
	metricsPort     = 8080
	grpcPort        = 6565
)

var (
	serviceName = common.GetEnv("OTEL_SERVICE_NAME", "SessionManagerFunctionGrpcPluginServerGoDocker")
	logLevelStr = common.GetEnv("LOG_LEVEL", logrus.InfoLevel.String())
)

func main() {
	go func() {
		runtime.SetBlockProfileRate(1)
		runtime.SetMutexProfileFraction(10)
	}()
	logrus.Infof("starting app server..")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logrusLevel, err := logrus.ParseLevel(logLevelStr)
	if err != nil {
		logrusLevel = logrus.InfoLevel
	}
	logrusLogger := logrus.New()
	logrusLogger.SetLevel(logrusLevel)

	loggingOptions := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall, logging.PayloadReceived, logging.PayloadSent),
		logging.WithFieldsFromContext(func(ctx context.Context) logging.Fields {
			if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
				return logging.Fields{"traceID", span.TraceID().String()}
			}

			return nil
		}),
		logging.WithLevels(logging.DefaultClientCodeToLevel),
		logging.WithDurationField(logging.DurationToDurationField),
	}

	srvMetrics := promgrpc.NewServerMetrics()
	unaryServerInterceptors := []grpc.UnaryServerInterceptor{
		srvMetrics.UnaryServerInterceptor(),
		logging.UnaryServerInterceptor(common.InterceptorLogger(logrusLogger), loggingOptions...),
	}
	streamServerInterceptors := []grpc.StreamServerInterceptor{
		srvMetrics.StreamServerInterceptor(),
		logging.StreamServerInterceptor(common.InterceptorLogger(logrusLogger), loggingOptions...),
	}

	// Preparing the IAM authorization
	var tokenRepo repository.TokenRepository = sdkAuth.DefaultTokenRepositoryImpl()
	var configRepo repository.ConfigRepository = sdkAuth.DefaultConfigRepositoryImpl()
	var refreshRepo repository.RefreshTokenRepository = &sdkAuth.RefreshTokenImpl{RefreshRate: 0.8, AutoRefresh: true}

	oauthService := iam.OAuth20Service{
		Client:                 factory.NewIamClient(configRepo),
		TokenRepository:        tokenRepo,
		RefreshTokenRepository: refreshRepo,
		ConfigRepository:       configRepo,
	}

	if strings.ToLower(common.GetEnv("PLUGIN_GRPC_SERVER_AUTH_ENABLED", "true")) == "true" {
		refreshInterval := common.GetEnvInt("REFRESH_INTERVAL", 600)
		common.Validator = common.NewTokenValidator(oauthService, time.Duration(refreshInterval)*time.Second, true)
		err := common.Validator.Initialize(ctx)
		if err != nil {
			logrus.Infof(err.Error())
		}

		unaryServerInterceptors = append(unaryServerInterceptors, common.UnaryAuthServerIntercept)
		streamServerInterceptors = append(streamServerInterceptors, common.StreamAuthServerIntercept)
		logrus.Infof("added auth interceptors")
	}

	gRPCServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(unaryServerInterceptors...),
		grpc.ChainStreamInterceptor(streamServerInterceptors...),
	)
	service := &server.SessionManager{
		UnimplementedSessionManagerServer: registered_v1.UnimplementedSessionManagerServer{},
	}
	registered_v1.RegisterSessionManagerServer(gRPCServer, service)

	// Enable gRPC Reflection
	reflection.Register(gRPCServer)
	logrus.Infof("gRPC reflection enabled")

	// Enable gRPC health check
	grpc_health_v1.RegisterHealthServer(gRPCServer, health.NewServer())

	// Register Prometheus Metrics
	srvMetrics.InitializeMetrics(gRPCServer)
	prometheusRegistry := prometheus.NewRegistry()
	prometheusRegistry.MustRegister(
		prometheusCollectors.NewGoCollector(),
		prometheusCollectors.NewProcessCollector(prometheusCollectors.ProcessCollectorOpts{}),
		srvMetrics,
	)

	go func() {
		http.Handle(metricsEndpoint, promhttp.HandlerFor(prometheusRegistry, promhttp.HandlerOpts{}))
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", metricsPort), nil))
	}()
	logrus.Infof("serving prometheus metrics at: (:%d%s)", metricsPort, metricsEndpoint)

	// Set Tracer Provider
	tracerProvider, err := common.NewTracerProvider(serviceName, environment, id)
	if err != nil {
		logrus.Fatalf("failed to create tracer provider: %v", err)
		return
	}

	otel.SetTracerProvider(tracerProvider)
	defer func(ctx context.Context) {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			logrus.Fatal(err)
		}
	}(ctx)
	logrus.Infof("set tracer provider: (name: %s environment: %s id: %d)", serviceName, environment, id)

	// Set Text Map Propagator
	b := b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader))
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			b,
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)
	logrus.Infof("set text map propagator")

	// Start gRPC Server
	logrus.Infof("starting gRPC server..")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		logrus.Fatalf("failed to listen to tcp:%d: %v", grpcPort, err)

		return
	}
	go func() {
		if err = gRPCServer.Serve(lis); err != nil {
			logrus.Fatalf("failed to run gRPC server: %v", err)

			return
		}
	}()
	logrus.Infof("gRPC server started")
	logrus.Infof("app server started")

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	logrus.Infof("signal received")
}
