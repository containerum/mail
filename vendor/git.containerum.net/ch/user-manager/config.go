package main

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func setupLogger() error {
	switch gin.Mode() {
	case gin.TestMode, gin.DebugMode:
	case gin.ReleaseMode:
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
	viper.SetDefault("log_level", logrus.InfoLevel)
	level := logrus.Level(viper.GetInt("log_level"))
	if level > logrus.DebugLevel || level < logrus.PanicLevel {
		return errors.New("invalid log level")
	}
	logrus.SetLevel(level)
	return nil
}

func getListenAddr() string {
	viper.SetDefault("listen_addr", ":8111")
	return viper.GetString("listen_addr")
}
