package component

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetComponent(t *testing.T) {
	mapper := GetComponentMappper()
	mapperClone := GetComponentMappper()

	// Test singleton pattern
	assert.Equal(t, &mapper, &mapperClone)

	// Test GetComponent
	components := GetComponents("pingcap", "tiflow", "[{\"name\":\"component/test\",\"color\":\"d1fad7\"},{\"name\":\"severity/minor\",\"color\":\"fbca04\"},{\"name\":\"bug-from-internal-test\",\"color\":\"1d76db\"}]")
	assert.Equal(t, Component("tiflow"), components[0])

	components = GetComponents("pingcap", "tiflow", "[{\"name\":\"area/ticdc\",\"color\":\"d1fad7\"},{\"name\":\"type/cherry-pick-for-release-5.4\",\"color\":\"9ad662\"}]")
	assert.Equal(t, TIFLOW_CDC, components[0])

	components = GetComponents("pingcap", "tiflow", "[{\"name\":\"area/dm\",\"color\":\"d1fad7\"},{\"name\":\"area/ticdc\",\"color\":\"9ad662\"}]")
	assert.Equal(t, TIFLOW_DM, components[0])
	assert.Equal(t, TIFLOW_CDC, components[1])
}
