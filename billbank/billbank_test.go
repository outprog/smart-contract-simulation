package billbank

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeposit(t *testing.T) {
	b := New()

	b.Deposit(100.0, "ETH", "alice")
	assert.Equal(t, 100.0, b.Pools["ETH"].Supply)
	assert.Equal(t, 100.0, b.Pools["ETH"].SupplyBill)
	assert.Equal(t, 100.0, b.AccountBills["alice"]["ETH"])

	b.Deposit(101.0, "ETH", "bob")
	assert.Equal(t, 201.0, b.Pools["ETH"].Supply)
	assert.Equal(t, 201.0, b.Pools["ETH"].SupplyBill)
	assert.Equal(t, 101.0, b.AccountBills["bob"]["ETH"])

	b.Deposit(101.0, "DAI", "alice")
	assert.Equal(t, 101.0, b.Pools["DAI"].Supply)
	assert.Equal(t, 101.0, b.Pools["DAI"].SupplyBill)
	assert.Equal(t, 101.0, b.AccountBills["alice"]["DAI"])

	b.Deposit(101.0, "DAI", "alice")
	assert.Equal(t, 202.0, b.Pools["DAI"].Supply)
	assert.Equal(t, 202.0, b.Pools["DAI"].SupplyBill)
	assert.Equal(t, 202.0, b.AccountBills["alice"]["DAI"])
}

func TestWithdraw(t *testing.T) {
	b := New()

	b.Deposit(100.0, "ETH", "alice")
	assert.Equal(t, 100.0, b.Pools["ETH"].Supply)
	assert.Equal(t, 100.0, b.Pools["ETH"].SupplyBill)

	_, err := b.Withdraw(1.0, "ETH", "bob")
	assert.EqualError(t, errors.New("user not had deposit. user: bob"), err.Error())

	_, err = b.Withdraw(101.0, "ETH", "alice")
	assert.EqualError(t, errors.New("not enough bill for withdraw. user: alice, acutal bill: 100"), err.Error())

	b.Withdraw(99.0, "ETH", "alice")
	assert.Equal(t, 1.0, b.Pools["ETH"].Supply)
	assert.Equal(t, 1.0, b.Pools["ETH"].SupplyBill)
	assert.Equal(t, 1.0, b.AccountBills["alice"]["ETH"])
}

func TestBorrow(t *testing.T) {
	b := New()

	err := b.Borrow(10.0, "ETH", "bob")
	assert.EqualError(t, errors.New("not enough token for borrow. amount: 10, cash: 0"), err.Error())

	b.Deposit(100.0, "ETH", "alice")

	err = b.Borrow(10.0, "ETH", "bob")
	assert.NoError(t, err)
	assert.Equal(t, 10.0, b.Pools["ETH"].Borrow)
}

func TestLiquidate(t *testing.T) {
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
	amount, _ := b.Withdraw(90.0, "ETH", "alice")
	assert.Equal(t, 90.81000000000002, amount)
	assert.Equal(t, 10.08999999999999, b.Pools["ETH"].Supply)
	assert.Equal(t, 10.0, b.AccountBills["alice"]["ETH"])
}
