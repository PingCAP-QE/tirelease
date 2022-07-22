package repository

import (
	"testing"
	"tirelease/commons/database"

	"github.com/stretchr/testify/assert"
)

func TestBatchSelectEmployeesByGhLogins(t *testing.T) {
	var config = generateConfig()
	database.Connect(config)

	githubLogins := []string{"MimeLyc", "VelocityLight"}
	employees, err := BatchSelectEmployeesByGhLogins(githubLogins)
	assert.Nil(t, err)
	assert.NotNil(t, employees)
}
