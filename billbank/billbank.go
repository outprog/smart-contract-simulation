package billbank

import (
	"fmt"
	"log"
)

type tSymbol = string

type TokenPool struct {
	SupplyBill float64
	Supply     float64
	Borrow     float64

	// last liquidate blockNumber
	borrowIndex uint64
}

// GetCash Cash = Supply - Borrow
func (t *TokenPool) GetCash() float64 {
	return t.Supply - t.Borrow
}

type Billbank struct {
	Pools map[tSymbol]TokenPool

	// BlockNumber simulate
	BlockNumber uint64
	// borrowRate every block
	borrowRate float64
}

func New() *Billbank {
	return &Billbank{
		Pools: map[tSymbol]TokenPool{
			"ETH": TokenPool{},
			"DAI": TokenPool{},
		},
		BlockNumber: 1,
		borrowRate:  0.01,
	}
}

func (b *Billbank) Deposit(amount float64, symbol string) error {
	b.liquidate(symbol)
	pool := b.getPool(symbol)

	// liquidate supply

	// update pool
	bill := amount
	if pool.SupplyBill != 0 && pool.Supply != 0 {
		bill = amount * (pool.SupplyBill / pool.Supply)
	}
	pool.SupplyBill += bill
	pool.Supply += amount
	b.Pools[symbol] = pool

	// send bill to user
	// sendBill(user, bill)

	return nil
}

func (b *Billbank) Withdraw(bill float64, symbol string) error {
	b.liquidate(symbol)
	pool := b.getPool(symbol)

	// liquidate supply

	// check balance of supply
	amount := bill * (pool.Supply / pool.SupplyBill)
	if bill > pool.SupplyBill || amount > pool.Supply {
		return fmt.Errorf("not enough token for withdraw. amount: %v, supply %v", amount, pool.Supply)
	}

	// update pool
	pool.SupplyBill -= bill
	pool.Supply -= amount
	b.Pools[symbol] = pool

	// send token to user
	// sendToken(user, amount)

	return nil
}

func (b *Billbank) Borrow(amount float64, symbol string) error {
	b.liquidate(symbol)
	pool := b.getPool(symbol)

	// check cash of pool
	if amount > pool.GetCash() {
		return fmt.Errorf("not enough token for borrow. amount: %v, cash: %v", amount, pool.Supply-pool.Borrow)
	}

	// update borrow
	pool.Borrow += amount
	b.Pools[symbol] = pool

	return nil
}

func (b *Billbank) Repay(amount float64, symbol string) error {
	b.liquidate(symbol)
	return nil
}

func (b *Billbank) liquidate(symbol string) {
	pool := b.getPool(symbol)

	receivable := 0.0
	if pool.Borrow != 0.0 {
		receivable = pool.Borrow * b.borrowRate * float64(b.BlockNumber-pool.borrowIndex)
	}
	pool.Supply += receivable

	// update pool
	pool.borrowIndex = b.BlockNumber
	b.Pools[symbol] = pool
}

func (b *Billbank) getPool(symbol string) (pool TokenPool) {
	var ok bool
	if pool, ok = b.Pools[symbol]; !ok {
		log.Panicf("not support token: %v", symbol)
	}
	return
}
