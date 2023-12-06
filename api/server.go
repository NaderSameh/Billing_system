package api

import (
	_ "github.com/naderSameh/billing_system/docs"
	"github.com/naderSameh/billing_system/limiter"
	"github.com/naderSameh/billing_system/worker"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

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
