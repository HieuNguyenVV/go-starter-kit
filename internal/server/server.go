package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-starter-kit/internal/log"
	"go-starter-kit/internal/pkg/database"
	"go-starter-kit/internal/server/config"
	"go-starter-kit/internal/server/middleware"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	logger     log.Logger
	config     *config.Config
	httpServer *gin.Engine
	postgres   *database.Postgres
}

func NewServer(config *config.Config,
	logger log.Logger,
	httpServer *gin.Engine,
	postgres *database.Postgres) *Server {

	{
		httpServer.GET("/healthz", func(c *gin.Context) {
			logger.Infof("Hello my friend....")
			logger.Infof("Hello")
			c.Status(http.StatusOK)
		})

		httpServer.GET("/readyz", func(c *gin.Context) {
			ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
			defer cancel()

			g, ctx := errgroup.WithContext(ctx)
			g.Go(postgres.Ping)

			if err := g.Wait(); err != nil {
				c.Status(http.StatusInternalServerError)
				return
			}
			c.Status(http.StatusOK)
		})
	}

	return &Server{
		logger:     logger,
		config:     config,
		httpServer: httpServer,
		postgres:   postgres,
	}
}

func (s *Server) Run() {
	sigint := make(chan os.Signal, 1)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%v", s.config.Server.Port),
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           s.httpServer,
	}
	go func() {
		s.logger.Infof("Server is running on port: %v", s.config.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("server is running error")
		}
	}()

	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	<-sigint
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		s.logger.Error("Server forced to shutdown")
	}
	s.postgres.Shutdown()

	s.logger.Info("Server exiting")
}

func NewHTTPServer(
	logger log.Logger,
	cfg *config.Config) *gin.Engine {
	var engine *gin.Engine
	if cfg.Gim.Debug {
		engine = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		engine = gin.New()
	}

	engine.Use(corsMiddleware, middleware.Cors(), middleware.Gzip(), middleware.Tx(logger))
	return engine
}

func corsMiddleware(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Next()
}
