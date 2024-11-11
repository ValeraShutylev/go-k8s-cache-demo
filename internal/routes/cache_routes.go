package routes

import (
	"cache-demo/internal/handlers"
	"cache-demo/internal/middleware"
	"cache-demo/internal/routes/health"
	"cache-demo/internal/routes/version"
	"cache-demo/internal/routes/metrics"
	"cache-demo/utils"

	"github.com/gin-gonic/gin"
)

const ( 
	ROOT_API = "/"
	CACHE_API = "/objects"
	PROBES_API = "/probes"
	AUDIT_ENABLED = "AUDIT_ENABLED"
)



type Router struct {
	GinRouter *gin.Engine
	cacheHandler *handlers.CacheHandler
}

func NewCacheRouter(cacheHandler *handlers.CacheHandler, state *utils.GlobalState) *Router {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	if utils.GetEnvAsBool(AUDIT_ENABLED, false) {
		router.Use(middleware.RequestLogger())
		router.Use(middleware.ResponseLogger())
	}
	router.Use(gin.Recovery())
	router.Use(func(c *gin.Context) {
		c.Set("state", state)
		c.Next()
	})

	r := &Router{GinRouter: router, cacheHandler: cacheHandler}

	r.registerCacheRoutes()
	r.registerHealthRoutes(state)
	r.registerVersionRoute() 
	r.registerMetricsRoute()
	return r
}

func(r *Router) registerCacheRoutes() {
	rg := r.GinRouter.Group(CACHE_API)
	rg.GET("/:objectId", r.cacheHandler.GetObjectById())
	rg.PUT("/:objectId", r.cacheHandler.PutObjectById())
}

func(r *Router) registerHealthRoutes(state *utils.GlobalState) {
	rg := r.GinRouter.Group(PROBES_API)
	health.HealthRoutes(rg, state)
}

func(r *Router) registerMetricsRoute() {
	rg := r.GinRouter.Group(ROOT_API)
	metrics.MetricsRoutes(rg)
}

func(r *Router) registerVersionRoute() {
	rg := r.GinRouter.Group(ROOT_API)
	version.AddAppVersionRoutes(rg)
	
}





