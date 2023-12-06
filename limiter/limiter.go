package limiter

import (
	"github.com/gin-gonic/gin"
	libredis "github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	limiter "github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
)

type Limiter interface {
	SetupRateLimiter() (gin.HandlerFunc, error)
}

type RateLimiterWithRedis struct {
	Middleware gin.HandlerFunc
	limit      string
}

func NewRateLimiter(limit string) (Limiter, error) {
	limiter := RateLimiterWithRedis{
		limit: limit,
	}
	return &limiter, nil
}

func (r *RateLimiterWithRedis) SetupRateLimiter() (gin.HandlerFunc, error) {

	rate, err := limiter.NewRateFromFormatted(r.limit)
	if err != nil {
		return nil, err
	}

	client := libredis.NewClient(&libredis.Options{Addr: viper.GetString("REDDIS_ADDR")})

	// Create a store with the redis client.
	store, err := sredis.NewStoreWithOptions(client, limiter.StoreOptions{
		Prefix:   "limiter_gin",
		MaxRetry: 3,
	})
	if err != nil {
		return nil, err
	}

	// Create a new middleware with the limiter instance.
	middleware := mgin.NewMiddleware(limiter.New(store, rate))

	return middleware, nil

}
