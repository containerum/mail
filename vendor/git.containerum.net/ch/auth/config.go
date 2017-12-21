package main

import (
	"encoding/base64"
	"errors"
	"io/ioutil"
	"strings"
	"time"

	"git.containerum.net/ch/auth/storages"
	"git.containerum.net/ch/auth/token"
	"git.containerum.net/ch/grpc-proto-files/auth"
	"github.com/dgrijalva/jwt-go"
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go-opentracing"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tidwall/buntdb"
)

func appendError(errs []string, err error) []string {
	if err != nil {
		return append(errs, err.Error())
	} else {
		return errs
	}
}

func setError(errs []string) error {
	if len(errs) == 0 {
		return nil
	} else {
		return errors.New(strings.Join(errs, "\n"))
	}
}

func getJWTConfig() (cfg token.JWTIssuerValidatorConfig, err error) {
	var errs []string

	cfg.SigningMethod = jwt.GetSigningMethod(viper.GetString("jwt_signing_method"))
	if cfg.SigningMethod == nil {
		errs = append(errs, "signing method not found")
	}

	cfg.Issuer = viper.GetString("issuer")

	viper.SetDefault("access_token_lifetime", 15*time.Minute)
	cfg.AccessTokenLifeTime = viper.GetDuration("access_token_lifetime")
	if cfg.AccessTokenLifeTime <= 0 {
		errs = append(errs, "access token lifetime is invalid or not set")
	}

	viper.SetDefault("refresh_token_lifetime", 48*time.Hour)
	cfg.RefreshTokenLifeTime = viper.GetDuration("refresh_token_lifetime")
	if cfg.RefreshTokenLifeTime <= cfg.AccessTokenLifeTime {
		errs = append(errs, "refresh token lifetime must be greater than access token lifetime")
	}

	signingKeyFile := viper.GetString("jwt_signing_key_file")
	validationKeyFile := viper.GetString("jwt_validation_key_file")
	signingKeyFileCont, err := ioutil.ReadFile(signingKeyFile)
	errs = appendError(errs, err)
	validationKeyFileCont, err := ioutil.ReadFile(validationKeyFile)
	errs = appendError(errs, err)
	switch cfg.SigningMethod.(type) {
	case *jwt.SigningMethodRSA, *jwt.SigningMethodRSAPSS:
		cfg.SigningKey, err = jwt.ParseRSAPrivateKeyFromPEM(signingKeyFileCont)
		errs = appendError(errs, err)
		cfg.ValidationKey, err = jwt.ParseRSAPublicKeyFromPEM(validationKeyFileCont)
		errs = appendError(errs, err)
	case *jwt.SigningMethodECDSA:
		cfg.SigningKey, err = jwt.ParseECPrivateKeyFromPEM(signingKeyFileCont)
		errs = appendError(errs, err)
		cfg.ValidationKey, err = jwt.ParseECPublicKeyFromPEM(validationKeyFileCont)
		errs = appendError(errs, err)
	default:
		signingKeyBuf := make([]byte, base64.StdEncoding.DecodedLen(len(signingKeyFileCont)))
		validationKeyBuf := make([]byte, base64.StdEncoding.DecodedLen(len(validationKeyFileCont)))
		_, err := base64.StdEncoding.Decode(signingKeyBuf, signingKeyFileCont)
		errs = appendError(errs, err)
		cfg.SigningKey = signingKeyBuf
		_, err = base64.StdEncoding.Decode(validationKeyBuf, validationKeyFileCont)
		errs = appendError(errs, err)
		cfg.ValidationKey = validationKeyBuf
	}

	err = setError(errs)
	return
}

func getTokenIssuerValidator() (iv token.IssuerValidator, err error) {
	viper.SetDefault("tokens", "jwt")
	tokens := viper.GetString("tokens")
	switch tokens {
	case "jwt":
		cfg, err := getJWTConfig()
		if err != nil {
			return nil, err
		}
		iv = token.NewJWTIssuerValidator(cfg)
		return iv, nil
	default:
		return nil, errors.New("invalid token issuer-validator")
	}
}

func getBuntDBStorageConfig(tokenFactory token.IssuerValidator) (cfg storages.BuntDBStorageConfig, err error) {
	var errs []string
	viper.SetDefault("storage_file", "storage.db")
	cfg.File = viper.GetString("storage_file")

	viper.SetDefault("bunt_syncpolicy", buntdb.EverySecond)
	cfg.BuntDBConfig.SyncPolicy = buntdb.SyncPolicy(viper.GetInt("bunt_synpolicy"))
	switch cfg.BuntDBConfig.SyncPolicy {
	case buntdb.EverySecond, buntdb.Never, buntdb.Always:
	default:
		errs = append(errs, "invalid bunt_syncpolicy")
	}

	viper.SetDefault("bunt_autoshrink_disabled", false)
	cfg.BuntDBConfig.AutoShrinkDisabled = viper.GetBool("bunt_autoshrink_disabled")

	if viper.IsSet("bunt_autoshrink_minsize") {
		cfg.BuntDBConfig.AutoShrinkMinSize = viper.GetInt("bunt_autoshrink_minsize")
	}

	if viper.IsSet("bunt_autoshrink_percentage") {
		cfg.BuntDBConfig.AutoShrinkPercentage = viper.GetInt("bunt_autoshrink_percentage")
	}

	cfg.TokenFactory = tokenFactory

	return cfg, setError(errs)
}

func getStorage() (storage auth.AuthServer, err error) {
	tokenFactory, err := getTokenIssuerValidator()
	if err != nil {
		return nil, err
	}

	viper.SetDefault("storage", "buntdb")
	switch viper.GetString("storage") {
	case "buntdb":
		var cfg storages.BuntDBStorageConfig
		cfg, err = getBuntDBStorageConfig(tokenFactory)
		if err != nil {
			return nil, err
		}
		storage, err = storages.NewBuntDBStorage(cfg)
		return
	default:
		return nil, errors.New("invalid storage")
	}
}

func getZipkinCollector() (collector zipkintracer.Collector, err error) {
	viper.SetDefault("zipkin_collector", "nop")
	switch viper.GetString("zipkin_collector") {
	case "http":
		collector, err = zipkintracer.NewHTTPCollector(viper.GetString("zipkin_http_collector_url"))
	case "kafka":
		collector, err = zipkintracer.NewKafkaCollector(viper.GetStringSlice("zipkin_kafka_collector_addrs"))
	case "scribe":
		collector, err = zipkintracer.NewScribeCollector(viper.GetString("zipkin_scribe_collector_addr"),
			viper.GetDuration("zipkin_scribe_collector_duration"))
	case "nop":
		collector = zipkintracer.NopCollector{}
	default:
		err = errors.New("invalid zipkin collector")
	}
	return
}

func getTracer(hostPort, service string) (tracer opentracing.Tracer, err error) {
	viper.SetDefault("tracer", "zipkin")
	switch viper.GetString("tracer") {
	case "zipkin":
		collector, err := getZipkinCollector()
		if err != nil {
			return nil, err
		}
		tracer, err = zipkintracer.NewTracer(zipkintracer.NewRecorder(collector,
			viper.GetBool("zipkin_recorder_debug"), hostPort, service))
	default:
		err = errors.New("invalid opentracing tracer found")
	}
	return
}

func logModeSetup() error {
	viper.SetDefault("log_mode", "text")
	switch viper.GetString("log_mode") {
	case "text":
		logrus.SetFormatter(&logrus.TextFormatter{})
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	default:
		return errors.New("invalid log mode")
	}
	return nil
}

func logLevelSetup() error {
	viper.SetDefault("log_level", logrus.InfoLevel)
	level := logrus.Level(viper.GetInt("log_level"))
	if level > logrus.DebugLevel || level < logrus.PanicLevel {
		return errors.New("invalid log level")
	}
	logrus.SetLevel(level)
	return nil
}
