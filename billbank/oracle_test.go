package billbank

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOracle(t *testing.T) {
	o := NewOracle()

	assert.Equal(t, 0.0, o.GetPrice("ETH"))

	o.SetPrice("ETH", 100.1)
	assert.Equal(t, 100.1, o.GetPrice("ETH"))
}
