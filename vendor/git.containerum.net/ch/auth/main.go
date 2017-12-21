package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

func logExit(err error) {
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}

func main() {
	viper.SetEnvPrefix("ch_auth")
	viper.AutomaticEnv()

	if err := logLevelSetup(); err != nil {
		fmt.Println(err)
		return
	}

	if err := logModeSetup(); err != nil {
		fmt.Println(err)
		return
	}

	viper.SetDefault("http_listenaddr", ":8080")
	httpTracer, err := getTracer(viper.GetString("http_listenaddr"), "ch-auth-rest")
	logExit(err)

	viper.SetDefault("grpc_listenaddr", ":8888")
	grpcTracer, err := getTracer(viper.GetString("grpc_listenaddr"), "ch-auth-grpc")
	logExit(err)

	storage, err := getStorage()
	logExit(err)

	RunServers(
		NewHTTPServer(viper.GetString("http_listenaddr"), httpTracer, storage),
		NewGRPCServer(viper.GetString("grpc_listenaddr"), grpcTracer, storage),
	)
}
