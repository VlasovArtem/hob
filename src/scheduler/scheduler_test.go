package scheduler

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func Test_Add(t *testing.T) {
	service := NewSchedulerService()

	channel := make(chan bool)

	itemId := uuid.New()
	entryId, err := service.Add(itemId, "* * * * *", func() {
		channel <- true
	})

	object := (service).(*SchedulerServiceObject)

	assert.Nil(t, err)
	assert.True(t, entryId > 0)

	go func() {
		object.cron.Entry(entryId).Job.Run()
	}()

	assert.True(t, <-channel)

	assert.NotEqual(t, cron.Entry{}, object.cron.Entry(entryId))
	id, ok := object.entries[itemId]

	assert.Equal(t, entryId, id)
	assert.True(t, ok)
}

func Test_Add_WithExistingScheduler(t *testing.T) {
	service := NewSchedulerService()

	itemId := uuid.New()
	entryId, err := service.Add(itemId, "* * * * *", func() { log.Println("Function") })

	object := (service).(*SchedulerServiceObject)

	assert.Nil(t, err)
	assert.True(t, entryId > 0)

	assert.NotEqual(t, cron.Entry{}, object.cron.Entry(entryId))
	id, ok := object.entries[itemId]

	assert.Equal(t, entryId, id)
	assert.True(t, ok)

	entryId, err = service.Add(itemId, "* * * * *", func() { log.Println("Function") })

	assert.Equal(t, cron.EntryID(0), entryId)
	assert.Equal(t, errors.New(fmt.Sprintf("scheduler for the entity id %s exists", itemId)), err)
}

func Test_Add_WithInvalidScheduler(t *testing.T) {
	service := NewSchedulerService()

	itemId := uuid.New()
	entryId, err := service.Add(itemId, "invalid", func() { log.Println("Function") })

	assert.NotNil(t, err)
	assert.Equal(t, cron.EntryID(0), entryId)

	object := (service).(*SchedulerServiceObject)

	id, ok := object.entries[itemId]

	assert.Equal(t, cron.EntryID(0), id)
	assert.False(t, ok)
}

func Test_Remove(t *testing.T) {
	service := NewSchedulerService()

	itemId := uuid.New()
	entryId, err := service.Add(itemId, "* * * * *", func() { log.Println("Function") })

	assert.Nil(t, err)
	assert.True(t, entryId > 0)

	err = service.Remove(itemId)

	assert.Nil(t, err)
}

func Test_Remove_WithNotExisingId(t *testing.T) {
	service := NewSchedulerService()

	itemId := uuid.New()

	err := service.Remove(itemId)

	assert.Equal(t, errors.New(fmt.Sprintf("Scheduler with is %s not found", itemId)), err)
}

func Test_Update(t *testing.T) {
	service := NewSchedulerService()
	itemId := uuid.New()

	entryId, err := service.Add(itemId, "* * * * *", func() {})

	object := (service).(*SchedulerServiceObject)

	assert.Nil(t, err)
	assert.True(t, entryId > 0)

	object.cron.Entry(entryId).Job.Run()

	assert.NotEqual(t, cron.Entry{}, object.cron.Entry(entryId))
	id, ok := object.entries[itemId]

	assert.Equal(t, entryId, id)
	assert.True(t, ok)

	channel := make(chan bool)
	entryId, err = service.Update(itemId, "* * * * *", func() {
		channel <- true
	})

	assert.Nil(t, err)
	assert.True(t, entryId > 0)

	go func() {
		object.cron.Entry(entryId).Job.Run()
	}()

	assert.True(t, <-channel)

	assert.NotEqual(t, cron.Entry{}, object.cron.Entry(entryId))
	id, ok = object.entries[itemId]

	assert.Equal(t, entryId, id)
	assert.True(t, ok)
}

func Test_Update_WithNotExists(t *testing.T) {
	service := NewSchedulerService()
	itemId := uuid.New()

	entryId, err := service.Update(itemId, "* * * * *", func() {})

	assert.Equal(t, cron.EntryID(0), entryId)
	assert.Equal(t, errors.New(fmt.Sprintf("scheduler for the entity id %s not exists", itemId)), err)
}
