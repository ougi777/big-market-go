package job

import "github.com/robfig/cron/v3"

type Scheduler struct {
	cron *cron.Cron
}

func NewScheduler() *Scheduler {
	return &Scheduler{cron: cron.New(cron.WithSeconds())}
}

func (s *Scheduler) Add(spec string, job func()) (cron.EntryID, error) {
	return s.cron.AddFunc(spec, job)
}

func (s *Scheduler) Start() {
	s.cron.Start()
}

func (s *Scheduler) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
}
