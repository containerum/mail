package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"

	"runtime/debug"

	"git.containerum.net/ch/auth/routes"
	"git.containerum.net/ch/grpc-proto-files/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/husobee/vestigo"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type HTTPServer struct {
	listenAddr string
	router     *vestigo.Router
}

func NewHTTPServer(listenAddr string, tracer opentracing.Tracer, storage auth.AuthServer) *HTTPServer {
	router := vestigo.NewRouter()
	routes.SetupRoutes(router, tracer, storage)
	return &HTTPServer{
		listenAddr: listenAddr,
		router:     router,
	}
}

func (s *HTTPServer) Run() error {
	logrus.WithField("listenAddr", s.listenAddr).Info("Starting HTTP server")
	return http.ListenAndServe(s.listenAddr, s.router)
}

type GRPCServer struct {
	listenAddr string
	server     *grpc.Server
}

func panicHandler(p interface{}) (err error) {
	logrus.Errorf("panic: %v", p)
	debug.PrintStack()
	return fmt.Errorf("panic: %v", p)
}

func NewGRPCServer(listenAddr string, tracer opentracing.Tracer, storage auth.AuthServer) *GRPCServer {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			otgrpc.OpenTracingServerInterceptor(tracer, otgrpc.LogPayloads()),
			grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(panicHandler)),
			grpc_logrus.UnaryServerInterceptor(logrus.WithField("component", "grpc_server")),
		)),
	)
	auth.RegisterAuthServer(server, storage)
	return &GRPCServer{
		listenAddr: listenAddr,
		server:     server,
	}
}

func (s *GRPCServer) Run() error {
	logrus.WithField("listenAddr", s.listenAddr).Infof("Starting GRPC server")
	listener, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	return s.server.Serve(listener)
}

type Runnable interface {
	Run() error
}

func RunServers(servers ...Runnable) {
	wg := &sync.WaitGroup{}
	wg.Add(len(servers))
	for _, server := range servers {
		go func(s Runnable) {
			if err := s.Run(); err != nil {
				logrus.Errorf("run server: %v", err)
				os.Exit(1)
			}
			wg.Done()
		}(server)
	}
	wg.Wait()
}
