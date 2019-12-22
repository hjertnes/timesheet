package database

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenTest(t *testing.T) {
	var database = OpenTest()
	assert.NotNil(t, database.Db)
}

func TestOpen(t *testing.T) {
	var database = Open()
	assert.NotNil(t, database.Db)
}

func TestInMemory(t *testing.T) {
	var database = OpenInMemory()
	assert.NotNil(t, database.Db)
}

func TestOpenFail(t *testing.T) {
	os.Setenv("HOME", "/")
	assert.Panics(t, func() { Open() })

}
