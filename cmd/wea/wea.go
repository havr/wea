package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/havr/wea/pkg/clients/openweather"
	"github.com/havr/wea/pkg/clients/wiki"
	"github.com/havr/wea/pkg/handler"
	"github.com/havr/wea/pkg/service"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

// Config defines external app configuration using environment variables.
type Config struct {
	ServeAt  string `envconfig:"SERVE_AT" required:"true"`
	OWAPIKey string `envconfig:"OW_API_KEY" required:"true"`
}

func main() {
	var config Config
	if err := envconfig.Process("WE", &config); err != nil {
		fmt.Println("config:", err)
		os.Exit(1)
	}

	// for demo purposes only
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println("logger:", err)
		os.Exit(1)
	}
	zap.ReplaceGlobals(logger)

	ow := openweather.NewClient(config.OWAPIKey)
	wiki := wiki.NewClient()
	wea := service.NewWea(ow, wiki)
	h := handler.New(wea)
	zap.L().Info("running at", zap.String("host", config.ServeAt))
	if err := http.ListenAndServe(config.ServeAt, h); err != nil && err != context.Canceled {
		zap.L().Error("server error", zap.Error(err))
	}
}
