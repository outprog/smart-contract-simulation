package billbank

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeposit(t *testing.T) {
	b := New()

	b.Deposit(100.0, "ETH")
	assert.Equal(t, 100.0, b.Pools["ETH"].Supply)
	assert.Equal(t, 100.0, b.Pools["ETH"].SupplyBill)

	b.Deposit(101.0, "ETH")
	assert.Equal(t, 201.0, b.Pools["ETH"].Supply)
	assert.Equal(t, 201.0, b.Pools["ETH"].SupplyBill)

	b.Deposit(101.0, "DAI")
	assert.Equal(t, 101.0, b.Pools["DAI"].Supply)
	assert.Equal(t, 101.0, b.Pools["DAI"].SupplyBill)
}

func TestWithdraw(t *testing.T) {
	b := New()

	b.Deposit(100.0, "ETH")
	assert.Equal(t, 100.0, b.Pools["ETH"].Supply)
	assert.Equal(t, 100.0, b.Pools["ETH"].SupplyBill)

	err := b.Withdraw(101.0, "ETH")
	assert.EqualError(t, errors.New("not enough token: 101 > 100"), err.Error())
	assert.Equal(t, 100.0, b.Pools["ETH"].Supply)
	assert.Equal(t, 100.0, b.Pools["ETH"].SupplyBill)

	b.Withdraw(99.0, "ETH")
	assert.Equal(t, 1.0, b.Pools["ETH"].Supply)
	assert.Equal(t, 1.0, b.Pools["ETH"].SupplyBill)
}
