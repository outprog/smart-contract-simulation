package money

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSupply(t *testing.T) {
	mc := New()

	user := "Alice"
	amount := 100.0
	symbol := "ETH"

	amount2 := 50.0
	symbol2 := "DAI"

	amount3 := 1.0
	symbol3 := "TFT"

	mc.Supply(amount, symbol, user)
	assert.Equal(t, amount, mc.Markets[symbol].supply)
	assert.Equal(t, amount, mc.SupplyBalances[user][symbol])

	mc.Supply(amount, symbol, user)
	assert.Equal(t, amount*2, mc.Markets[symbol].supply)
	assert.Equal(t, amount*2, mc.SupplyBalances[user][symbol])

	mc.Supply(amount2, symbol2, user)
	assert.Equal(t, amount2, mc.Markets[symbol2].supply)
	assert.Equal(t, amount2, mc.SupplyBalances[user][symbol2])

	err := mc.Supply(amount3, symbol3, user)
	assert.EqualError(t, errors.New("not support token: TFT"), err.Error())
}

func TestBorrow(t *testing.T) {
	mc := New()
	user := "Alice"

	amount := 100.0
	symbol := "ETH"

	user2 := "Bob"
	amount2 := 33.3

	// borrow with err
	err := mc.Borrow(amount, "TFT", user)
	assert.EqualError(t, errors.New("not support token: TFT"), err.Error())
	err = mc.Borrow(amount, symbol, user)
	assert.EqualError(t, errors.New("not enough cash: 0"), err.Error())

	// borrow success
	mc.Supply(amount, symbol, user)
	err = mc.Borrow(amount2, symbol, user2)
	assert.NoError(t, err)
	assert.Equal(t, amount2, mc.Markets[symbol].borrow)
	assert.Equal(t, amount2, mc.BorrowBalances[user2][symbol])

	// borrow again
	err = mc.Borrow(amount2, symbol, user2)
	assert.NoError(t, err)
	assert.Equal(t, amount2*2, mc.Markets[symbol].borrow)
	assert.Equal(t, amount2*2, mc.BorrowBalances[user2][symbol])
}
