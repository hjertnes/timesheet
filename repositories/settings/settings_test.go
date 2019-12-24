package settings

import (
	"testing"

	"github.com/hjertnes/timesheet/database"
	"github.com/hjertnes/timesheet/utils"
	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	var database = database.OpenInMemory()

	var repo = New(database.Db)

	err := repo.AddOrUpdate("Test", "Test")

	utils.ErrorHandler(err)

	var items, _ = repo.GetAll()

	assert.Len(t, items, 1)

	err = repo.AddOrUpdate("Test", "Test2")

	utils.ErrorHandler(err)

	items, _ = repo.GetAll()

	assert.Len(t, items, 1)

	err = repo.AddOrUpdate("Test2", "Test2")

	utils.ErrorHandler(err)

	items, _ = repo.GetAll()
	assert.Len(t, items, 2)
}

func TestList(t *testing.T) {
	var database = database.OpenInMemory()

	var repo = New(database.Db)

	err := repo.AddOrUpdate("Test", "Test")

	utils.ErrorHandler(err)

	err = repo.AddOrUpdate("Test2", "Test")

	utils.ErrorHandler(err)

	var items, _ = repo.GetAll()

	assert.Len(t, items, 2)
}

func TestGetOne(t *testing.T) {
	var database = database.OpenInMemory()

	var repo = New(database.Db)

	err := repo.AddOrUpdate("Test", "Test")

	utils.ErrorHandler(err)

	err = repo.AddOrUpdate("Test2", "Test")

	utils.ErrorHandler(err)

	_, err = repo.GetOne("Test")
	assert.Nil(t, err)
}

func TestDeleteAll(t *testing.T) {
	var database = database.OpenInMemory()

	var repo = New(database.Db)

	err := repo.AddOrUpdate("Test", "Test")

	utils.ErrorHandler(err)

	err = repo.AddOrUpdate("Test2", "Test")

	utils.ErrorHandler(err)

	err = repo.DeleteAll()

	utils.ErrorHandler(err)

	assert.Nil(t, err)
}

func TestExist(t *testing.T) {
	var database = database.OpenInMemory()

	var repo = New(database.Db)

	err := repo.AddOrUpdate("Test", "Test")

	utils.ErrorHandler(err)

	var ex1, _ = repo.Exist("Test")

	var ex2, _ = repo.Exist("Test2")

	assert.True(t, ex1)

	assert.False(t, ex2)
}
