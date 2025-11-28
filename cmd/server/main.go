package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"git.uhomes.net/uhs-go/go-bisub/internal/config"
	fxmodules "git.uhomes.net/uhs-go/go-bisub/internal/pkg/fx"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fxmodules.ConfigModule,
		fxmodules.LoggerModule,
		fxmodules.DatabaseModule,
		fxmodules.RedisModule,
		fxmodules.RepositoryModule,
		fxmodules.ServiceModule,
		fxmodules.HandlerModule,
		fxmodules.MiddlewareModule,
		fxmodules.HTTPModule,
		fx.Invoke(startServer),
	)

	app.Run()
}

func startServer(lc fx.Lifecycle, engine *gin.Engine, cfg *config.Config) {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      engine,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				slog.Info("Server starting", "port", cfg.Server.Port)
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					slog.Error("Failed to start server", "error", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			slog.Info("Shutting down server...")
			ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()
			return srv.Shutdown(ctx)
		},
	})
}


