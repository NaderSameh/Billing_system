package cron

import (
	"time"

	"github.com/go-co-op/gocron"
	db "github.com/naderSameh/billing_system/db/sqlc"
	"github.com/rs/zerolog/log"
)

const (
	CronSecond     = "Second"
	CronMin        = "Min"
	CronHour       = "Hour"
	CronMonth      = "Month"
	CronEndOfMonth = "EndOfMonth"
)

type CronJobScheduler interface {
	UpdateMonthlyBilling(store db.Store) error
	StartCron() error
	StopCron() error
}

type GoCronScheduler struct {
	S *gocron.Scheduler
}

type CronJobConfiguration struct {
	Interval int
	Unit     string
	Name     string
}

type CronJob struct {
	job *gocron.Job
}

func NewGoCronScheduler(t *time.Location, config CronJobConfiguration, store db.Store) (CronJobScheduler, CronJob) {
	s := gocron.NewScheduler(t)
	scheduler := GoCronScheduler{
		S: s,
	}
	var job *gocron.Job
	J := CronJob{
		job: job,
	}

	var err error
	if config.Unit == CronEndOfMonth {
		J.job, err = s.Every(config.Interval).MonthLastDay().Name(config.Name).Do(func() error { return scheduler.UpdateMonthlyBilling(store) })
	} else if config.Unit == CronSecond {
		J.job, err = s.Every(config.Interval).Second().Name(config.Name).Do(func() error { return scheduler.UpdateMonthlyBilling(store) })
	}
	if err != nil {
		panic("failed to schedule task")
	}

	J.job.RegisterEventListeners(
		gocron.WhenJobReturnsError(func(jobName string, err error) {
			log.Error().Err(err).Str("CronName", jobName).Msg("cron task failed")
		}),
	)

	return &scheduler, J
}
func (cron *GoCronScheduler) StartCron() error {
	cron.S.StartAsync()
	return nil
}

func (cron *GoCronScheduler) StopCron() error {
	cron.S.Stop()
	return nil
}
