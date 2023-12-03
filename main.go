package main

import (
	"database/sql"
	"runtime"

	_ "github.com/lib/pq"
	"github.com/naderSameh/billing_system/api"
	db "github.com/naderSameh/billing_system/db/sqlc"
	"github.com/naderSameh/billing_system/util"

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
// @schemes	http
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

	server.Start(viper.GetString("SERVER_ADDRESS"))
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start server:")
	}

}
