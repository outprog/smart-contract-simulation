package billbank

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLiquidateTokenPool(t *testing.T) {
	b := New()
	b.borrowRate = 0.01

	b.BlockNumber = 1
	b.Deposit(100.0, "ETH", "alice")
	b.Borrow(10.0, "ETH", "bob")

	b.BlockNumber = 2
	b.Deposit(0.0, "ETH", "alice")
	assert.Equal(t, 100.1, b.Pools["ETH"].Supply)

	b.BlockNumber = 4
	b.Borrow(10.0, "ETH", "bob")
	assert.Equal(t, 100.3, b.Pools["ETH"].Supply)

	b.BlockNumber = 10
	b.Borrow(0.0, "ETH", "bob")
	assert.Equal(t, 100.3+20*6*0.01, b.Pools["ETH"].Supply)
}

func TestDepositInterest(t *testing.T) {
	b := New()

	b.BlockNumber = 1
	b.Deposit(100.0, "ETH", "alice")
	b.Borrow(10.0, "ETH", "bob")

	b.BlockNumber = 10
	// this 90.0 is bill. bill price auto increase when borrow is not 0.
	// amount, _ := b.Withdraw(90.0, "ETH", "alice")
	// assert.Equal(t, 90.81000000000002, amount)
	// assert.Equal(t, 10.08999999999999, b.Pools["ETH"].Supply)
	// assert.Equal(t, 10.0, b.AccountDepositBills["alice"]["ETH"])
}
