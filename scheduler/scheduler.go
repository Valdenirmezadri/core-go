package scheduler

import (
	"github.com/robfig/cron/v3"
)

type EntryID int

type SpecType string

func (SpecType) New(spec string) SpecType {
	return SpecType(spec)
}

func (s SpecType) String() string {
	return string(s)
}

const (
	MidNight    SpecType = "00 00 * * *"
	EveryMinute SpecType = "1 * * * *"
	EveryHour   SpecType = "0 * * * *"
	TwoAM       SpecType = "0 2 * * *"
	FiveAM      SpecType = "0 5 * * *"
)

type Scheduler interface {
	AddFunc(spec SpecType, cmd func()) (EntryID, error)
	Remove(ID EntryID)
}

type scheduler struct {
	cron *cron.Cron
}

func New() Scheduler {
	s := &scheduler{
		cron: cron.New(),
	}

	s.cron.Start()

	return s
}

// AddFunc adds a func to the Cron to be run on the given schedule. The spec is parsed using the time zone of this Cron instance as the default. An opaque ID is returned that can be used to later remove it.
func (s *scheduler) AddFunc(spec SpecType, cmd func()) (EntryID, error) {
	id, err := s.cron.AddFunc(spec.String(), cmd)
	if err != nil {
		return 0, err

	}

	return EntryID(id), nil
}

func (s *scheduler) Remove(ID EntryID) {
	s.cron.Remove(cron.EntryID(ID))
}
