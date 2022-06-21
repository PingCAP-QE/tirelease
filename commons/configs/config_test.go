package configs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	t.Skip()
	LoadConfig("../../config.yaml")
	assert.Equal(t, true, Config.Github.AccessToken != "")
}
