package version

import (
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
)

func AddAppVersionRoutes(rg *gin.RouterGroup) {
	rg.GET("/version", appVersion)
}

func appVersion(ctx *gin.Context) {
	appVersion := os.Getenv("APP_VERSION")
	if appVersion == "" {
		appVersion = "v1.0.0"
	}

	ctx.String(
		http.StatusOK,
		appVersion,
	)
}