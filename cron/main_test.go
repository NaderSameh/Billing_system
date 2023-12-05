package cron

import (
	"os"
	"testing"
	"time"

	db "github.com/naderSameh/billing_system/db/sqlc"
	"github.com/rs/zerolog/log"
)

func NewTestCron(t *testing.T, store db.Store, config CronJobConfiguration) (CronJobScheduler, CronJob) {

	scheduler, job := NewGoCronScheduler(time.UTC, config, store)

	err := scheduler.StartCron()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start scheduler")
	}
	return scheduler, job
}
func TestMain(m *testing.M) {

	os.Exit(m.Run())
}
