package task

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTask(t *testing.T) {
	task := NewRefreshEmployeeTask()
	assert.Equal(t, reflect.TypeOf(task.ITaskExecution), reflect.TypeOf(task))
}
