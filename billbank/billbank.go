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
	BorrowBill float64
	Borrow     float64

	// last liquidate blockNumber
	borrowIndex uint64
}

// GetCash Cash = Supply - Borrow
func (t *TokenPool) GetCash() float64 {
	return t.Supply - t.Borrow
}

type Billbank struct {
	// internal account for token bill(deposit)
	AccountDepositBills map[tUser]map[tSymbol]tBill
	// internal account for token bill(borrow)
	AccountBorrowBills map[tUser]map[tSymbol]tBill

	Pools map[tSymbol]TokenPool

	// BlockNumber simulate
	BlockNumber uint64
	// borrowRate every block
	borrowRate float64
}

func New() *Billbank {
	return &Billbank{
		AccountDepositBills: map[tUser]map[tSymbol]tBill{},
		AccountBorrowBills:  map[tUser]map[tSymbol]tBill{},
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
	if accountBill, ok := b.AccountDepositBills[user]; ok {
		if _, ok := accountBill[symbol]; ok {
			b.AccountDepositBills[user][symbol] += bill
		} else {
			b.AccountDepositBills[user][symbol] = bill
		}
	} else {
		b.AccountDepositBills[user] = map[tSymbol]tBill{symbol: bill}
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
	if _, ok := b.AccountDepositBills[user]; !ok {
		return 0, fmt.Errorf("user not had deposit. user: %v", user)
	}
	if _, ok := b.AccountDepositBills[user][symbol]; !ok {
		return 0, fmt.Errorf("user not had deposit. user: %v, token: %v", user, symbol)
	}
	if bill > b.AccountDepositBills[user][symbol] {
		return 0, fmt.Errorf("not enough bill for withdraw. user: %v, acutal bill: %v", user, b.AccountDepositBills[user][symbol])
	}
	// check balance of supply
	if amount > pool.GetCash() {
		return 0, fmt.Errorf("not enough token for withdraw. amount: %v, cash %v", amount, pool.GetCash())
	}

	// calcuate amount
	amount = bill * (pool.Supply / pool.SupplyBill)

	// update user account bill
	b.AccountDepositBills[user][symbol] -= bill

	// update pool
	pool.SupplyBill -= bill
	pool.Supply -= amount
	b.Pools[symbol] = pool

	return
}

func (b *Billbank) Borrow(amount float64, symbol, user string) error {
	b.liquidate(symbol)
	pool := b.getPool(symbol)

	// check cash of pool
	if amount > pool.GetCash() {
		return fmt.Errorf("not enough token for borrow. amount: %v, cash: %v", amount, pool.GetCash())
	}

	// TODO
	// calcuate bill
	// bill := amount
	// if pool.BorrowBill != 0 && pool.Borrow != 0 {
	// 	bill = amount * (pool.BorrowBill / pool.Borrow)
	// }

	// update user account borrow
	if accountBorrow, ok := b.AccountBorrowBills[user]; ok {
		if _, ok := accountBorrow[symbol]; ok {
			b.AccountBorrowBills[user][symbol] += amount
		} else {
			b.AccountBorrowBills[user][symbol] = amount
		}
	} else {
		b.AccountBorrowBills[user] = map[tSymbol]tBill{symbol: amount}
	}

	// update borrow
	// pool.BorrowBill += bill
	pool.Borrow += amount
	b.Pools[symbol] = pool

	return nil
}

func (b *Billbank) Repay(amount float64, symbol, user string) error {
	b.liquidate(symbol)
	pool := b.getPool(symbol)

	// check borrow
	if _, ok := b.AccountBorrowBills[user]; !ok {
		return fmt.Errorf("user not had borrow. user: %v", user)
	}
	if _, ok := b.AccountBorrowBills[user][symbol]; !ok {
		return fmt.Errorf("user not had borrow. user: %v, token: %v", user, symbol)
	}
	if amount > b.AccountBorrowBills[user][symbol] {
		return fmt.Errorf("too much amount to repay. user: %v, need repay: %v", user, b.AccountBorrowBills[user][symbol])
	}

	// update user account borrow
	b.AccountBorrowBills[user][symbol] -= amount

	// update borrow
	pool.Borrow -= amount
	b.Pools[symbol] = pool

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
