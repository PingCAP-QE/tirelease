package repository

import (
	"testing"
	"tirelease/commons/configs"

	"github.com/stretchr/testify/assert"
)

func TestSelectAllHrEmployees(t *testing.T) {
	configs.LoadConfig("../../config.yaml")
	InitHrEmployeeDB()
	hrEmployees, err := SelectAllHrEmployee()
	assert.Nil(t, err)
	assert.NotNil(t, hrEmployees)
}
