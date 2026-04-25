package providers

import (
	"net/http"
	appconfig "transport/internal/infrastructure/config/app"

	"github.com/gin-gonic/gin"
)

func HTTP(cfg appconfig.Config, routers *gin.Engine) *http.Server {
	return &http.Server{
		Addr:    cfg.HTTP.Port,
		Handler: routers,
	}
}
