package billbank

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeposit(t *testing.T) {
	b := New()

	b.Deposit(100.0, "ETH")
	assert.Equal(t, 100.0, b.Pools["ETH"].Supply)

	b.Deposit(101.0, "ETH")
	assert.Equal(t, 201.0, b.Pools["ETH"].Supply)

	b.Deposit(101.0, "DAI")
	assert.Equal(t, 101.0, b.Pools["DAI"].Supply)
}
