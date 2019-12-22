package event

import (
	"testing"
	"time"

	"github.com/hjertnes/timesheet/database"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	var database = database.OpenInMemory()
	var repo = NewEventRepository(database.Db)
	repo.Add(time.Now(), time.Now(), false, false)
	var all, _ = repo.GetAll()
	assert.Len(t, all, 1)
}

func TestList(t *testing.T) {
	var database = database.OpenInMemory()
	var repo = NewEventRepository(database.Db)
	repo.Add(time.Now(), time.Now(), false, false)
	repo.Add(time.Now(), time.Now(), false, false)
	var all, _ = repo.GetAll()
	assert.Len(t, all, 2)
}

func TestDeleteAll(t *testing.T) {
	var database = database.OpenInMemory()
	var repo = NewEventRepository(database.Db)
	repo.Add(time.Now(), time.Now(), false, false)
	repo.Add(time.Now(), time.Now(), false, false)
	var err = repo.DeleteAll()
	var all, _ = repo.GetAll()
	assert.Nil(t, err)
	assert.Len(t, all, 0)
}

func TestDeleteOne(t *testing.T) {
	var database = database.OpenInMemory()
	var repo = NewEventRepository(database.Db)
	repo.Add(time.Now(), time.Now(), false, false)
	repo.Add(time.Now(), time.Now(), false, false)
	var all, _ = repo.GetAll()
	var first = all[0]
	var err = repo.Delete(int(first.ID))
	assert.Nil(t, err)
	all, _ = repo.GetAll()
	assert.Len(t, all, 1)
}
func TestDeleteOne_Error(t *testing.T) {
	var database = database.OpenInMemory()
	var repo = NewEventRepository(database.Db)
	repo.Add(time.Now(), time.Now(), false, false)
	repo.Add(time.Now(), time.Now(), false, false)
	var err = repo.Delete(123)
	assert.NotNil(t, err)
	all, _ := repo.GetAll()
	assert.Len(t, all, 2)
}
