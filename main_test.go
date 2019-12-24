package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	t.Run("Run", func(t *testing.T) {
		main()
		assert.Nil(t, nil)
	})
}
