package billbank

import "fmt"

type tSymbol = string

type TokenPool struct {
	SupplyBill float64
	Supply     float64
	Borrow     float64
}

type Billbank struct {
	Pools map[tSymbol]TokenPool

	// BlockNumber simulate
	BlockNumber uint64
}

func New() *Billbank {
	return &Billbank{
		Pools: map[tSymbol]TokenPool{
			"ETH": TokenPool{},
			"DAI": TokenPool{},
		},
		BlockNumber: 1,
	}
}

func (b *Billbank) Deposit(amount float64, symbol string) error {
	// check pool
	var pool TokenPool
	var ok bool
	if pool, ok = b.Pools[symbol]; !ok {
		return fmt.Errorf("not support token: %v", symbol)
	}

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
	// check pool
	var pool TokenPool
	var ok bool
	if pool, ok = b.Pools[symbol]; !ok {
		return fmt.Errorf("not support token: %v", symbol)
	}

	// liquidate supply

	// check balance of supply
	amount := bill * (pool.Supply / pool.SupplyBill)
	if bill > pool.SupplyBill || amount > pool.Supply {
		return fmt.Errorf("not enough token: %v > %v", amount, pool.Supply)
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
	return nil
}

func (b *Billbank) Repay(amount float64, symbol string) error {
	return nil
}
