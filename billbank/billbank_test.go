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
	assert.EqualError(t, errors.New("not enough token for withdraw. amount: 101, supply 100"), err.Error())
	assert.Equal(t, 100.0, b.Pools["ETH"].Supply)
	assert.Equal(t, 100.0, b.Pools["ETH"].SupplyBill)

	b.Withdraw(99.0, "ETH")
	assert.Equal(t, 1.0, b.Pools["ETH"].Supply)
	assert.Equal(t, 1.0, b.Pools["ETH"].SupplyBill)
}

func TestBorrow(t *testing.T) {
	b := New()

	err := b.Borrow(10.0, "ETH")
	assert.EqualError(t, errors.New("not enough token for borrow. amount: 10, cash: 0"), err.Error())

	b.Deposit(100.0, "ETH")

	err = b.Borrow(10.0, "ETH")
	assert.NoError(t, err)
	assert.Equal(t, 10.0, b.Pools["ETH"].Borrow)
}

func TestLiquidate(t *testing.T) {
	b := New()

	b.BlockNumber = 1
	b.Deposit(100.0, "ETH")
	b.Borrow(10.0, "ETH")

	b.BlockNumber = 2
	b.Deposit(0.0, "ETH")
	assert.Equal(t, 100.1, b.Pools["ETH"].Supply)

	b.BlockNumber = 4
	b.Borrow(10.0, "ETH")
	assert.Equal(t, 100.3, b.Pools["ETH"].Supply)

	b.BlockNumber = 10
	b.Borrow(0.0, "ETH")
	assert.Equal(t, 100.3+20*6*0.01, b.Pools["ETH"].Supply)
}
