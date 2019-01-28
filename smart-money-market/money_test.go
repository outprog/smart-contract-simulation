package money

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSupply(t *testing.T) {
	mc := New()
	user := "Alice"
	amount := 100.0
	symbol := "ETH"

	mc.Supply(amount, symbol, user)
	assert.Equal(t, amount, mc.Markets[symbol].supply)
	assert.Equal(t, amount, mc.SupplyBalances[user][symbol])

	mc.Supply(amount, symbol, user)
	assert.Equal(t, amount*2, mc.Markets[symbol].supply)
	assert.Equal(t, amount*2, mc.SupplyBalances[user][symbol])

	amount2 := 50.0
	symbol2 := "DAI"
	mc.Supply(amount2, symbol2, user)
	assert.Equal(t, amount2, mc.Markets[symbol2].supply)
	assert.Equal(t, amount2, mc.SupplyBalances[user][symbol2])
}

func TestBorrow(t *testing.T) {
	mc := New()
	user := "Alice"
	amount := 100.0
	symbol := "ETH"
	mc.Borrow(amount, "TFT", user)
	mc.Borrow(amount, symbol, user)
}
