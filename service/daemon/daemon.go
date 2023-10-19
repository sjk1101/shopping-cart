package daemon

import (
	"context"
	"sync"

	"shopping-cart/service/binder"
	"shopping-cart/service/daemon/schedule"

	"go.uber.org/dig"
)

var (
	once   sync.Once
	daemon *Daemon
)

func InitDaemon(in Daemon) {
	once.Do(func() {
		daemon = &in
	},
	)
}

type Daemon struct {
	dig.In

	ScheduleJob schedule.ScheduleInterface
}

func Run() {
	b := binder.New()
	if err := b.Provide(schedule.NewSchedule); err != nil {
		panic(err)
	}
	if err := b.Invoke(InitDaemon); err != nil {
		panic(err)
	}

	if daemon == nil {
		panic("daemon is nil")
	}

	ctx := context.Background()
	go daemon.ScheduleJob.Run(ctx)
}
