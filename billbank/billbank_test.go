package billbank

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLiquidateTokenPool(t *testing.T) {
	b := New()
	b.borrowRate = 0.01

	b.BlockNumber = 1
	b.Deposit(100.0, "ETH", "alice")

	b.BlockNumber = 3
	b.liquidate("ETH")
	assert.Equal(t, 100.0, b.Pools["ETH"].Supply)

	b.BlockNumber = 4
	b.Borrow(10.0, "ETH", "bob")

	b.BlockNumber = 10
	b.liquidate("ETH")
	assert.Equal(t, 10*math.Pow(1.01, 10-4), b.Pools["ETH"].Borrow)
	assert.Equal(t, 100.0+(10*math.Pow(1.01, 10-4)-10), b.Pools["ETH"].Supply)
}

func TestGrowth(t *testing.T) {
	b := New()
	b.borrowRate = 0.02

	b.BlockNumber = 1
	b.Deposit(100.0, "ETH", "alice")
	b.Borrow(10.0, "ETH", "bob")

	b.BlockNumber = 10
	growth := 10.0*math.Pow(1.02, 10-1) - 10.0
	assert.Equal(t, 10.0+growth, b.BorrowBalanceOf("ETH", "bob"))
	assert.Equal(t, 10.0, b.Pools["ETH"].Borrow)
	// assert.Equal(t, 100.0+growth, b.SupplyBalanceOf("ETH", "alice"))

	// after liquidate, pool changed
	b.liquidate("ETH")
	assert.Equal(t, 10.0+growth, b.BorrowBalanceOf("ETH", "bob"))
	assert.Equal(t, 10.0+growth, b.Pools["ETH"].Borrow)
}
