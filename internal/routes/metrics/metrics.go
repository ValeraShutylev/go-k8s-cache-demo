package metrics

import (
	"cache-demo/internal/repository"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var MapLengthGauge  = prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "memory_cache_items_count",
    Help: "The current length of the map",
})

func init() {
    prometheus.MustRegister(MapLengthGauge)
}

func MetricsRoutes(rg *gin.RouterGroup) {
	rg.GET("/metrics", gin.WrapH(promhttp.Handler()))
}

func UpdateMapLength(length int) {
    MapLengthGauge.Set(float64(length))
}

func ScrapCacheLenghtMetrics(mem *repository.MemoryCache) {
    for {
        UpdateMapLength(mem.Len())
        time.Sleep(10 * time.Second)
    }
}