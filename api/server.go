package api

import (
	"time"

	ginzap "github.com/gin-contrib/zap"
	_ "github.com/naderSameh/billing_system/docs"
	"github.com/naderSameh/billing_system/limiter"
	"github.com/naderSameh/billing_system/worker"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	db "github.com/naderSameh/billing_system/db/sqlc"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

var taskDistributor worker.TaskDistributor

func NewServer(store db.Store) (*Server, error) {
	server := &Server{
		store: store,
	}
	server.setupRouter()
	taskDistributor = worker.NewRedisDistributor(viper.GetString("REDDIS_ADDR"))

	return server, nil
}

func (server *Server) setupRouter() {
	r := gin.Default()

	SetupCORS(r)
	SetupRateLimiter(r)
	SetupLogger(r)
	SetupMetrics(r)

	r.POST("/batches", server.createBatch)
	r.DELETE("/batches/:batch_id", server.deleteBatch)
	r.GET("/batches", server.listBatches)
	r.PUT("/batches/:batch_id", server.updateBatch)

	r.POST("/bundles", server.createBundle)
	r.POST("/bundles/assign", server.assignBundleToCustomer)
	r.GET("/bundles", server.getBundles)
	r.DELETE("/bundles/:bundle_id", server.deleteBundle)

	r.POST("/orders", server.createOrder)
	r.PUT("/orders/:order_id", server.updateOrder)

	r.PUT("/payments_logs/:log_id", server.updatePaymentLog)
	r.DELETE("/payments_logs/:log_id", server.deletePayment)
	r.GET("/payments_logs", server.listPaymentLogs)

	r.GET("/charges", server.listChargesPerCustomer)

	//Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server.router = r
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func SetupCORS(router *gin.Engine) {
	router.Use(CORSMiddleware())
}

func SetupRateLimiter(router *gin.Engine) {
	limiter, err := limiter.NewRateLimiter(viper.GetString("RATE_LIMIT"))
	if err != nil {
		log.Error().Err(err).Msg("failed to setup rate limiter")
	}
	limiterMiddleware, err := limiter.SetupRateLimiter()
	router.Use(limiterMiddleware)
}

func SetupLogger(router *gin.Engine) {

	logger, _ := zap.NewProduction()
	// Add a ginzap middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	//   - RFC3339 with UTC time format.

	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))

	// Logs all panic to error log
	router.Use(ginzap.RecoveryWithZap(logger, true))

}
func SetupMetrics(router *gin.Engine) {

	// get global Monitor object
	m := ginmetrics.GetMonitor()

	// +optional set metric path, default /debug/metrics
	m.SetMetricPath("/metrics")
	// +optional set slow time, default 5s
	m.SetSlowTime(10)
	// +optional set request duration, default {0.1, 0.3, 1.2, 5, 10}
	// used to p95, p99
	m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})

	// set middleware for gin
	m.Use(router)

}
