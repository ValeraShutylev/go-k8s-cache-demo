package health

import (
	"cache-demo/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)



func HealthRoutes(rg *gin.RouterGroup, state *utils.GlobalState) {
	rg.GET("/health", health)
	rg.GET("/ready", ready)
}

func health(ctx *gin.Context) {
	ctx.String(
		http.StatusOK,
		"Ok")
}

func ready(ctx *gin.Context) {
	stateFromCtx, exists := ctx.Get("state")
	if exists {
		state := stateFromCtx.(*utils.GlobalState)
		if state.GetState() != utils.DONE {
			ctx.String(
				http.StatusServiceUnavailable,
				"Current state of Data transfer from storage to memory cache: ", state.GetState())
				return
		} 
	}
	ctx.String(
		http.StatusOK,
		"Ok")
}



