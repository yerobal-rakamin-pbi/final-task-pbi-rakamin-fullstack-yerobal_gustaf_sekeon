package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"rakamin-final-task/config"
	"rakamin-final-task/database"
	swagger "rakamin-final-task/docs"
	"rakamin-final-task/helpers/log"
	"rakamin-final-task/helpers/response"
	"rakamin-final-task/helpers/errors"
	"rakamin-final-task/middlewares"

	"context"
	"fmt"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type router struct {
	config      config.Application
	http        *gin.Engine
	db          *database.DB
	log         log.LogInterface
	response    response.Interface
	middlewares middlewares.Interface
}

func (r router) setupSwagger() {
	swagger.SwaggerInfo.Host = fmt.Sprintf("%s:%s", r.config.Server.Host, r.config.Server.Port)
	swagger.SwaggerInfo.Schemes = []string{"https", "http"}
}

var once = sync.Once{}

func Init(
	config config.Application,
	log log.LogInterface,
	db *database.DB,
) router {
	r := router{}

	// Initialize server with graceful shutdown
	once.Do(func() {
		gin.SetMode(gin.ReleaseMode)

		r.config = config
		r.http = gin.New()
		r.log = log
		r.db = db
		r.response = response.Init(log)
		r.middlewares = middlewares.Init(config, r.http, r.response)

		r.setupSwagger()
		r.RegisterMiddlewaresAndRoutes()
	})

	return r
}

func (r router) RegisterMiddlewaresAndRoutes() {
	// Global middleware
	r.http.Use(r.middlewares.SetCors())
	r.http.Use(r.middlewares.SetTimeout)
	r.http.Use(r.middlewares.AddFieldsToCtx)

	r.setupSwagger()

	// Global routes
	r.http.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.http.GET("/ping", r.ping)


	// 404 handler
	r.http.NoRoute(r.notFoundHandler)
}

func (r router) ping(c *gin.Context) {
	r.response.Success(c, "PONG!!!", nil, nil)
}

func (r router) notFoundHandler(c *gin.Context) {
	r.response.Error(c, errors.NotFound("Endpoint not found"))
}

func (r router) Run() {
	/*
		Create context that listens for the interrupt signal from the OS.
		This will allow us to gracefully shutdown the server.
	*/
	c := context.Background()
	ctx, stop := signal.NotifyContext(c, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	port := fmt.Sprintf(":8080")
	if r.config.Server.Port != "" {
		port = ":" + r.config.Server.Port
	}
	server := &http.Server{
		Addr:              port,
		Handler:           r.http,
		ReadHeaderTimeout: 2 * time.Second,
	}

	// Run the server in a goroutine so that it doesn't block the graceful shutdown handling below

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			r.log.Error(ctx, err.Error())
		}
	}()

	r.log.Info(context.Background(), "Server is running on port "+r.config.Server.Port)

	// Block until we receive our signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	r.log.Info(context.Background(), "Shutting down server...")

	// Create a deadline to wait for.
	quitCtx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	if err := server.Shutdown(quitCtx); err != nil {
		r.log.Fatal(quitCtx, fmt.Sprintf("Server Shutdown error: %s", err.Error()))
	}

	r.log.Info(context.Background(), "Server gracefully stopped")
}
