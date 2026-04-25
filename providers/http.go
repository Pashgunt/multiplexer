package providers

import (
	"net/http"
	"time"
	appconfig "transport/internal/infrastructure/config/app"

	"github.com/gin-gonic/gin"
)

func HTTP(cfg appconfig.Config, routers *gin.Engine) *http.Server {
	return &http.Server{
		Addr:              cfg.HTTP.Port,
		Handler:           routers,
		ReadTimeout:       time.Duration(cfg.HTTP.ReadTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(cfg.HTTP.ReadHeaderTimeout) * time.Second,
		WriteTimeout:      time.Duration(cfg.HTTP.WriteTimeout) * time.Second,
	}
}
