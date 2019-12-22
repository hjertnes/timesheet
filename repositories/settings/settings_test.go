package settings

import (
	"testing"

	"github.com/hjertnes/timesheet/database"
	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	var database = database.OpenInMemory()
	var repo = NewSettingsRepository(database.Db)
	repo.AddOrUpdate("Test", "Test")
	var items, _ = repo.GetAll()
	assert.Len(t, items, 1)
	repo.AddOrUpdate("Test", "Test2")
	items, _ = repo.GetAll()
	assert.Len(t, items, 1)
	repo.AddOrUpdate("Test2", "Test2")
	items, _ = repo.GetAll()
	assert.Len(t, items, 2)
}

func TestList(t *testing.T) {
	var database = database.OpenInMemory()
	var repo = NewSettingsRepository(database.Db)
	repo.AddOrUpdate("Test", "Test")
	repo.AddOrUpdate("Test2", "Test")
	var items, _ = repo.GetAll()
	assert.Len(t, items, 2)
}

func TestGetOne(t *testing.T) {
	var database = database.OpenInMemory()
	var repo = NewSettingsRepository(database.Db)
	repo.AddOrUpdate("Test", "Test")
	repo.AddOrUpdate("Test2", "Test")
	var _, err = repo.GetOne("Test")
	assert.Nil(t, err)

}

func TestDeleteAll(t *testing.T) {
	var database = database.OpenInMemory()
	var repo = NewSettingsRepository(database.Db)
	repo.AddOrUpdate("Test", "Test")
	repo.AddOrUpdate("Test2", "Test")
	var err = repo.DeleteAll()
	assert.Nil(t, err)
}

func TestExist(t *testing.T) {
	var database = database.OpenInMemory()
	var repo = NewSettingsRepository(database.Db)
	repo.AddOrUpdate("Test", "Test")
	var ex1, _ = repo.Exist("Test")
	var ex2, _ = repo.Exist("Test2")
	assert.True(t, ex1)
	assert.False(t, ex2)
}
