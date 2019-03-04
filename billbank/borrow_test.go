package billbank

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBorrow(t *testing.T) {
	b := New()

	err := b.Borrow(10.0, "ETH", "bob")
	assert.EqualError(t, err, "not enough token for borrow. amount: 10, cash: 0")

	b.Deposit(100.0, "ETH", "alice")

	err = b.Borrow(10.0, "ETH", "bob")
	assert.NoError(t, err)
	assert.Equal(t, 10.0, b.Pools["ETH"].Borrow)
	assert.Equal(t, 10.0, b.AccountBorrowBills["bob"]["ETH"])
	assert.Equal(t, 10.0, b.BorrowBalanceOf("ETH", "bob"))

	err = b.Borrow(13.0, "ETH", "bob")
	assert.NoError(t, err)
	assert.Equal(t, 23.0, b.Pools["ETH"].Borrow)
	assert.Equal(t, 23.0, b.AccountBorrowBills["bob"]["ETH"])
	assert.Equal(t, 23.0, b.BorrowBalanceOf("ETH", "bob"))
}

func TestRepay(t *testing.T) {
	b := New()

	err := b.Repay(10.0, "ETH", "bob")
	assert.EqualError(t, err, "too much amount to repay. user: bob, need repay: 0")

	b.Deposit(20.0, "ETH", "alice")
	b.Borrow(10.0, "ETH", "bob")

	err = b.Repay(11.0, "ETH", "bob")
	assert.EqualError(t, err, "too much amount to repay. user: bob, need repay: 10")

	b.Repay(4.0, "ETH", "bob")
	assert.Equal(t, 6.0, b.AccountBorrowBills["bob"]["ETH"])
	assert.Equal(t, 6.0, b.Pools["ETH"].Borrow)
	assert.Equal(t, 6.0, b.BorrowBalanceOf("ETH", "bob"))

	b.Repay(6.0, "ETH", "bob")
	assert.Equal(t, 0.0, b.AccountBorrowBills["bob"]["ETH"])
	assert.Equal(t, 0.0, b.Pools["ETH"].Borrow)
	assert.Equal(t, 0.0, b.BorrowBalanceOf("ETH", "bob"))
}

func TestBorrowValueOf(t *testing.T) {
	b := New()

	b.Deposit(10.0, "ETH", "alice")
	b.Borrow(10.0, "ETH", "bob")

	assert.Equal(t, 0.0, b.BorrowValueOf("ETH", "bob"))

	b.Oralcer.SetPrice("ETH", 100.1)
	assert.Equal(t, 1001.0, b.BorrowValueOf("ETH", "bob"))
	assert.Equal(t, 0.0, b.BorrowValueOf("ETH", "alice"))

	b.Oralcer.SetPrice("ETH", 100.2)
	assert.Equal(t, 1002.0, b.BorrowValueOf("ETH", "bob"))
}
