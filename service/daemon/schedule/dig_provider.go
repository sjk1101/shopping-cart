package schedule

import (
	"context"
	"sync"

	"shopping-cart/service/repository"

	"github.com/robfig/cron"
	"go.uber.org/dig"
)

var (
	self *schedule
	once sync.Once
)

type ScheduleInterface interface {
	Run(ctx context.Context)
}

func NewSchedule(in scheduleIn) ScheduleInterface {
	once.Do(func() {
		self = &schedule{
			in:      in,
			cronJob: cron.New(),
		}
	})

	return self
}

type scheduleIn struct {
	dig.In

	ProductRepo repository.ProductRepositoryInterface
}

type schedule struct {
	in scheduleIn

	cronJob *cron.Cron
}

type JobInterface interface {
	GetExpired() string
	Exec(ctx context.Context)
}

func (s *schedule) Run(ctx context.Context) {
	in := self.in

	// 註冊cron
	s.register(ctx, newMonitorCronJob(in))
	s.cronJob.Start()
}

func (s *schedule) register(ctx context.Context, cronJob JobInterface) {
	if err := s.cronJob.AddFunc(cronJob.GetExpired(), func() { cronJob.Exec(ctx) }); err != nil {
		panic(err)
	}
}
