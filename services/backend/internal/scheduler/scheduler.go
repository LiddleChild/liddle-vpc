package scheduler

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/LiddleChild/liddle-vpc/backend/internal/timeutil"
	"github.com/robfig/cron/v3"
)

type Task interface {
	Name() string
	Schedule() string
	Run(ctx context.Context) error
}

type Scheduler struct {
	cron *cron.Cron

	ctx    context.Context
	cancel context.CancelFunc
}

func New() *Scheduler {
	cron := cron.New(cron.WithLocation(timeutil.LocationBangkok))
	cron.Start()

	ctx, cancel := context.WithCancel(context.Background())

	return &Scheduler{
		cron:   cron,
		ctx:    ctx,
		cancel: cancel,
	}
}

// ScheduleTask schedules a cron task. Panic if task has invalid cron schedule.
func (s *Scheduler) ScheduleTask(task Task) {
	if _, err := s.cron.AddFunc(task.Schedule(), func() {
		if err := task.Run(s.ctx); err != nil {
			slog.Error("task failed", "name", task.Name(), "error", err)
			return
		}
	}); err != nil {
		panic(fmt.Errorf("failed to schedule task %s: %w", task.Name(), err))
	}
}

func (s *Scheduler) Stop(ctx context.Context) {
	s.cancel()

	select {
	case <-s.cron.Stop().Done():
	case <-ctx.Done():
	}
}
