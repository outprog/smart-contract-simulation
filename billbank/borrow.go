package billbank

import "fmt"

func (b *Billbank) BorrowBalanceOf(symbol, user string) float64 {
	pool := b.getPool(symbol)

	// check bill
	if _, ok := b.AccountBorrowBills[user]; !ok {
		return 0.0
	}
	bill := 0.0
	if b, ok := b.AccountBorrowBills[user][symbol]; ok {
		bill = b
	}
	if bill == 0.0 {
		return 0.0
	}

	// calcuate amount
	return bill * (pool.Borrow / pool.BorrowBill)
}

func (b *Billbank) Borrow(amount float64, symbol, user string) error {
	b.liquidate(symbol)
	pool := b.getPool(symbol)

	// check cash of pool
	if amount > pool.GetCash() {
		return fmt.Errorf("not enough token for borrow. amount: %v, cash: %v", amount, pool.GetCash())
	}

	// calcuate bill
	bill := amount
	if pool.BorrowBill != 0 && pool.Borrow != 0 {
		bill = amount * (pool.BorrowBill / pool.Borrow)
	}

	// update user account bill
	if accountBorrow, ok := b.AccountBorrowBills[user]; ok {
		if _, ok := accountBorrow[symbol]; ok {
			b.AccountBorrowBills[user][symbol] += bill
		} else {
			b.AccountBorrowBills[user][symbol] = bill
		}
	} else {
		b.AccountBorrowBills[user] = map[tSymbol]tBill{symbol: bill}
	}

	// update borrow
	pool.BorrowBill += bill
	pool.Borrow += amount
	b.Pools[symbol] = pool

	return nil
}

func (b *Billbank) Repay(amount float64, symbol, user string) error {
	b.liquidate(symbol)
	pool := b.getPool(symbol)

	// check borrow
	accountAmount := b.BorrowBalanceOf(symbol, user)
	if amount > accountAmount {
		return fmt.Errorf("too much amount to repay. user: %v, need repay: %v", user, accountAmount)
	}

	// calculate bill
	bill := amount * (pool.BorrowBill / pool.Borrow)

	// update user account borrow
	b.AccountBorrowBills[user][symbol] -= bill

	// update borrow
	pool.BorrowBill -= bill
	pool.Borrow -= amount
	b.Pools[symbol] = pool

	return nil
}
