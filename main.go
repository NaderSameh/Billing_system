package main

import (
	"database/sql"
	"runtime"
	"time"

	_ "github.com/lib/pq"
	"github.com/naderSameh/billing_system/api"
	"github.com/naderSameh/billing_system/cron"
	db "github.com/naderSameh/billing_system/db/sqlc"
	"github.com/naderSameh/billing_system/mail"
	"github.com/naderSameh/billing_system/util"
	"github.com/naderSameh/billing_system/worker"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

//	@title			Gin Swagger Example API
//	@version		1.0
//	@description	Ticketing support microservice
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Cypodsolutions
//	@contact.url	http://www.cypod.com/
//	@contact.email	naders@cypodsolutions.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost:8080
// @BasePath	/
// @schemes	http https
func main() {

	err := util.Loadconfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	conn, err := sql.Open(viper.GetString("DB_DRIVER"), viper.GetString("DB_SOURCE"))
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}
	conn.SetMaxOpenConns(runtime.NumCPU())
	n := conn.Stats().MaxOpenConnections
	log.Info().Int("connections", n).Msg("db max number of connections")

	store := db.NewStore(conn)

	server, err := api.NewServer(store)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	config := cron.CronJobConfiguration{
		Interval: 1,
		Unit:     cron.CronEndOfMonth,
		Name:     "Monthly Billing",
	}
	cron, _ := cron.NewGoCronScheduler(time.UTC, config, store)
	err = cron.StartCron()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start scheduler")
	}

	go runTaskProcessor(viper.GetString("REDDIS_ADDR"))

	server.Start(viper.GetString("SERVER_ADDRESS"))
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start server:")
	}

}

func runTaskProcessor(RedisAddress string) {
	mailer := mail.NewGmailSender(viper.GetString("GMAIL_NAME"), viper.GetString("GMAIL_EMAIL"), viper.GetString("GMAIL_PASS"))
	taskProcessor := worker.NewRedisTaskProcessor(RedisAddress, mailer)
	log.Info().Msg("start task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}
}
