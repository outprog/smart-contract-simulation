package billbank

import (
	"fmt"
	"log"
)

type tUser = string
type tSymbol = string
type tBill = float64

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
	AccountBills map[tUser]map[tSymbol]tBill

	Pools map[tSymbol]TokenPool

	// BlockNumber simulate
	BlockNumber uint64
	// borrowRate every block
	borrowRate float64
}

func New() *Billbank {
	return &Billbank{
		AccountBills: map[tUser]map[tSymbol]tBill{},
		Pools: map[tSymbol]TokenPool{
			"ETH": TokenPool{},
			"DAI": TokenPool{},
		},
		BlockNumber: 1,
		borrowRate:  0.01,
	}
}

func (b *Billbank) Deposit(amount float64, symbol, user string) error {
	b.liquidate(symbol)
	pool := b.getPool(symbol)

	// calcuate bill
	bill := amount
	if pool.SupplyBill != 0 && pool.Supply != 0 {
		bill = amount * (pool.SupplyBill / pool.Supply)
	}

	// update user account bill
	if accountBill, ok := b.AccountBills[user]; ok {
		if _, ok := accountBill[symbol]; ok {
			b.AccountBills[user][symbol] += bill
		} else {
			b.AccountBills[user][symbol] = bill
		}
	} else {
		b.AccountBills[user] = map[tSymbol]tBill{symbol: bill}
	}

	// update pool
	pool.SupplyBill += bill
	pool.Supply += amount
	b.Pools[symbol] = pool

	return nil
}

func (b *Billbank) Withdraw(bill float64, symbol, user string) (amount float64, err error) {
	b.liquidate(symbol)
	pool := b.getPool(symbol)

	// check bill
	if _, ok := b.AccountBills[user]; !ok {
		return 0, fmt.Errorf("user not had deposit. user: %v", user)
	}
	if _, ok := b.AccountBills[user][symbol]; !ok {
		return 0, fmt.Errorf("user not had deposit. user: %v, token: %v", user, symbol)
	}
	if bill > b.AccountBills[user][symbol] {
		return 0, fmt.Errorf("not enough bill for withdraw. user: %v, acutal bill: %v", user, b.AccountBills[user][symbol])
	}
	// check balance of supply
	if amount > pool.GetCash() {
		return 0, fmt.Errorf("not enough token for withdraw. amount: %v, cash %v", amount, pool.GetCash())
	}

	// calcuate amount
	amount = bill * (pool.Supply / pool.SupplyBill)

	// update user account bill
	b.AccountBills[user][symbol] -= bill

	// update pool
	pool.SupplyBill -= bill
	pool.Supply -= amount
	b.Pools[symbol] = pool

	return
}

func (b *Billbank) Borrow(amount float64, symbol string) error {
	b.liquidate(symbol)
	pool := b.getPool(symbol)

	// check cash of pool
	if amount > pool.GetCash() {
		return fmt.Errorf("not enough token for borrow. amount: %v, cash: %v", amount, pool.GetCash())
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
