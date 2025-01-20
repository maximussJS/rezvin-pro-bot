package http

import (
	"context"
	"crypto/tls"
	"fmt"
	"go.uber.org/dig"
	"net/http"
	"rezvin-pro-bot/config"
	"rezvin-pro-bot/constants"
	"rezvin-pro-bot/internal/logger"
	"rezvin-pro-bot/router"
	"sync"
	"time"
)

type IHttpServer interface {
	Start()
}

type httpServerDependencies struct {
	dig.In

	ShutdownWaitGroup *sync.WaitGroup `name:"ShutdownWaitGroup"`
	ShutdownContext   context.Context `name:"ShutdownContext"`
	Router            router.IRouter  `name:"Router"`
	Logger            logger.ILogger  `name:"Logger"`
	Config            config.IConfig  `name:"Config"`
}

type httpServer struct {
	shutdownWaitGroup *sync.WaitGroup
	shutdownContext   context.Context
	logger            logger.ILogger
	config            config.IConfig
	router            router.IRouter
	server            http.Server
}

func NewHttpServer(deps httpServerDependencies) *httpServer {
	tlsConfig := &tls.Config{
		ClientAuth: tls.NoClientCert,
		MinVersion: tls.VersionTLS11,
	}

	return &httpServer{
		shutdownWaitGroup: deps.ShutdownWaitGroup,
		shutdownContext:   deps.ShutdownContext,
		logger:            deps.Logger,
		config:            deps.Config,
		router:            deps.Router,
		server: http.Server{
			Addr:         deps.Config.HttpPort(),
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  30 * time.Second,
			TLSConfig:    tlsConfig,
			Handler:      deps.Router.GetHttpRouter(),
		},
	}
}

func (s *httpServer) Shutdown(ctx context.Context) error {
	s.logger.Log("Calling shutdown on http server")

	return s.server.Shutdown(ctx)
}

func (s *httpServer) Start() {
	defer s.shutdownWaitGroup.Done()

	s.shutdownWaitGroup.Add(1)

	go func() {
		if s.config.AppEnv() == constants.DevelopmentEnv {
			s.logger.Log(fmt.Sprintf("Starting http server on %s", s.server.Addr))

			if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				s.logger.Fatal(fmt.Sprintf("Failed to start http server: %s", err))
			}
		} else {
			panic("https server not implemented")
			//server.logger.Log(fmt.Sprintf("Starting https server on port %s", port))
			//if err := server.server.ListenAndServeTLS(server.config.SSLCertPath(), server.config.SSLKeyPath()); err != nil && err != http.ErrServerClosed {
			//	server.logger.Fatal(fmt.Sprintf("Failed to start https server: %s", err))
			//}
		}
	}()

	select {
	case <-s.shutdownContext.Done():
		s.logger.Log("Shutting down server gracefully...")
		shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelShutdown()

		err := s.Shutdown(shutdownCtx)
		if err != nil {
			s.logger.Error(fmt.Sprintf("Failed to shutdown server gracefully: %s", err))
		}
	}

	s.logger.Log("Server stopped")
}
