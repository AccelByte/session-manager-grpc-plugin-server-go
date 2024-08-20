package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"

	sessionmanager "accelbyte.net/session-manager-grpc-plugin-server-go/pkg/pb"
	"accelbyte.net/session-manager-grpc-plugin-server-go/pkg/session"

	"accelbyte.net/session-manager-grpc-plugin-server-go/pkg/config"
	"accelbyte.net/session-manager-grpc-plugin-server-go/pkg/utils"
	"accelbyte.net/session-manager-grpc-plugin-server-go/pkg/utils/envelope"
	"github.com/caarlos0/env"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	//nolint:gochecknoglobals
	buildDate = "unknown"

	//nolint:gochecknoglobals
	revisionID = "unknown"

	//nolint:gochecknoglobals
	gitHash = "unknown"

	//nolint:gochecknoglobals
	rolesSeedingVersion = "unknown"
)

const serviceName = "session-manager-grpc-plugin-server-go"

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, os.Kill) //nolint:staticcheck

	rootCtx, cancel := context.WithCancel(context.Background())
	scope := envelope.NewRootScope(rootCtx, serviceName, utils.MakeTraceID(serviceName))
	defer scope.Finish()

	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	logrus.Info("Session manager Grpc Plugin")
	logrus.Infof("RevisionID: %s, Build Date: %s, Git Hash: %s Roles Seeding Version: %s\n", revisionID, buildDate, gitHash, rolesSeedingVersion)

	cfg := &config.Config{}

	flag.Usage = func() {
		flag.CommandLine.SetOutput(os.Stdout)
		for _, val := range cfg.HelpDocs() {
			//nolint:forbidigo
			fmt.Println(val)
		}

		//nolint:forbidigo
		fmt.Println("")
		flag.PrintDefaults()
	}
	flag.Parse()

	err := env.Parse(cfg)
	if err != nil {
		logrus.Error("unable to parse environment variables: ", err)

		return
	}

	grpcServer := grpc.NewServer()
	service := &session.SessionManager{
		UnimplementedSessionManagerServer: sessionmanager.UnimplementedSessionManagerServer{},
	}

	sessionmanager.RegisterSessionManagerServer(grpcServer, service)

	// Enable gRPC Reflection
	reflection.Register(grpcServer)
	logrus.Infof("gRPC reflection enabled")

	gRPCListener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", cfg.GRPCPort))
	if err != nil {
		logrus.Fatalf("unable to create gRPC listener: %v", err)
	}

	go func() {
		logrus.Infof("Serving gRPC, listens at 0.0.0.0:%d", cfg.GRPCPort)
		if errServeGRPC := grpcServer.Serve(gRPCListener); err != nil {
			logrus.Fatalf("failed to serve gRPC: %v", errServeGRPC)
		}
	}()

	//nolint:gosimple
	select {
	case <-sigCh:
		cancel()
	}
}