package billbank

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeposit(t *testing.T) {
	b := New()

	b.Deposit(100.0, "ETH", "alice")
	assert.Equal(t, 100.0, b.Pools["ETH"].Supply)
	assert.Equal(t, 100.0, b.Pools["ETH"].SupplyBill)
	assert.Equal(t, 100.0, b.AccountDepositBills["alice"]["ETH"])
	assert.Equal(t, 100.0, b.SupplyBalanceOf("ETH", "alice"))

	b.Deposit(101.0, "ETH", "bob")
	assert.Equal(t, 201.0, b.Pools["ETH"].Supply)
	assert.Equal(t, 201.0, b.Pools["ETH"].SupplyBill)
	assert.Equal(t, 101.0, b.AccountDepositBills["bob"]["ETH"])
	assert.Equal(t, 101.0, b.SupplyBalanceOf("ETH", "bob"))

	b.Deposit(101.0, "DAI", "alice")
	assert.Equal(t, 101.0, b.Pools["DAI"].Supply)
	assert.Equal(t, 101.0, b.Pools["DAI"].SupplyBill)
	assert.Equal(t, 101.0, b.AccountDepositBills["alice"]["DAI"])
	assert.Equal(t, 101.0, b.SupplyBalanceOf("DAI", "alice"))

	b.Deposit(101.0, "DAI", "alice")
	assert.Equal(t, 202.0, b.Pools["DAI"].Supply)
	assert.Equal(t, 202.0, b.Pools["DAI"].SupplyBill)
	assert.Equal(t, 202.0, b.AccountDepositBills["alice"]["DAI"])
	assert.Equal(t, 202.0, b.SupplyBalanceOf("DAI", "alice"))
}

func TestWithdraw(t *testing.T) {
	b := New()

	err := b.Withdraw(1.0, "ETH", "bob")
	assert.EqualError(t, err, "not enough amount for withdraw. user: bob, acutal amount: 0")

	b.Deposit(100.0, "ETH", "alice")
	assert.Equal(t, 100.0, b.Pools["ETH"].Supply)
	assert.Equal(t, 100.0, b.Pools["ETH"].SupplyBill)
	assert.Equal(t, 100.0, b.SupplyBalanceOf("ETH", "alice"))

	err = b.Withdraw(101.0, "ETH", "alice")
	assert.EqualError(t, err, "not enough amount for withdraw. user: alice, acutal amount: 100")

	b.Withdraw(99.0, "ETH", "alice")
	assert.Equal(t, 1.0, b.Pools["ETH"].Supply)
	assert.Equal(t, 1.0, b.Pools["ETH"].SupplyBill)
	assert.Equal(t, 1.0, b.AccountDepositBills["alice"]["ETH"])
	assert.Equal(t, 1.0, b.SupplyBalanceOf("ETH", "alice"))
}

func TestSupplyValueOf(t *testing.T) {
	b := New()

	b.Deposit(100.1, "ETH", "alice")
	assert.Equal(t, 0.0, b.SupplyValueOf("ETH", "alice"))

	b.Oralcer.SetPrice("ETH", 100.0)
	assert.Equal(t, 10010.0, b.SupplyValueOf("ETH", "alice"))
}
