package scheduler

import (
	"context"
	"errors"
	"fmt"
	"github.com/VlasovArtem/hob/src/common/dependency"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

type SchedulingSpecification string

const (
	HOURLY   SchedulingSpecification = "@hourly"
	DAILY    SchedulingSpecification = "@daily"
	WEEKLY   SchedulingSpecification = "@weekly"
	MONTHLY  SchedulingSpecification = "@monthly"
	ANNUALLY SchedulingSpecification = "@annually"
)

type SchedulerServiceObject struct {
	cron    *cron.Cron
	entries map[uuid.UUID]cron.EntryID
}

func NewSchedulerService() ServiceScheduler {
	schedulerService := &SchedulerServiceObject{
		cron:    cron.New(),
		entries: make(map[uuid.UUID]cron.EntryID),
	}
	schedulerService.cron.Start()

	return schedulerService
}

func (s *SchedulerServiceObject) Initialize(factory dependency.DependenciesFactory) interface{} {
	return factory.Add(NewSchedulerService())
}

type ServiceScheduler interface {
	Add(scheduledItemId uuid.UUID, scheduleSpec string, scheduleFunc func()) (cron.EntryID, error)
	Remove(id uuid.UUID) error
	Stop() context.Context
	Update(id uuid.UUID, scheduleSpec string, scheduleFunc func()) (cron.EntryID, error)
}

func (s *SchedulerServiceObject) Add(scheduledItemId uuid.UUID, scheduleSpec string, scheduleFunc func()) (entryID cron.EntryID, err error) {
	if _, ok := s.entries[scheduledItemId]; ok {
		return entryID, errors.New(fmt.Sprintf("scheduler for the entity id %s exists", scheduledItemId))
	}
	if entryID, err = s.cron.AddFunc(scheduleSpec, scheduleFunc); err != nil {
		return entryID, err
	} else {
		s.entries[scheduledItemId] = entryID

		return entryID, err
	}
}

func (s *SchedulerServiceObject) Remove(scheduledItemId uuid.UUID) error {
	if entryId, ok := s.entries[scheduledItemId]; !ok {
		return errors.New(fmt.Sprintf("Scheduler with is %s not found", scheduledItemId))
	} else {
		s.cron.Remove(entryId)
	}
	return nil
}

func (s *SchedulerServiceObject) Stop() context.Context {
	return s.cron.Stop()
}

func (s *SchedulerServiceObject) Update(id uuid.UUID, scheduleSpec string, scheduleFunc func()) (entryID cron.EntryID, err error) {
	if entryID, ok := s.entries[id]; !ok {
		return entryID, errors.New(fmt.Sprintf("scheduler for the entity id %s not exists", id))
	} else {
		s.cron.Remove(entryID)
		if entryID, err = s.cron.AddFunc(scheduleSpec, scheduleFunc); err != nil {
			return entryID, err
		} else {
			s.entries[id] = entryID
			return entryID, err
		}
	}
}
