package billbank

import "fmt"

func (b *Billbank) SupplyBalanceOf(symbol, user string) float64 {
	pool := b.getPool(symbol)

	// check bill
	if _, ok := b.AccountDepositBills[user]; !ok {
		return 0.0
	}
	bill := 0.0
	if b, ok := b.AccountDepositBills[user][symbol]; ok {
		bill = b
	}
	if bill == 0.0 {
		return 0.0
	}

	// calcuate amount
	// current block liquidated, growth is zero
	growth := b.calculateGrowth(symbol)
	return bill * ((pool.Supply + growth) / pool.SupplyBill)
}

func (b *Billbank) SupplyValueOf(symbol, user string) float64 {
	return b.SupplyBalanceOf(symbol, user) * b.Oralcer.GetPrice(symbol)
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

func (b *Billbank) Withdraw(amount float64, symbol, user string) (err error) {
	b.liquidate(symbol)
	pool := b.getPool(symbol)

	// check account balance
	accountAmount := b.SupplyBalanceOf(symbol, user)
	if amount > accountAmount {
		return fmt.Errorf("not enough amount for withdraw. user: %v, acutal amount: %v", user, accountAmount)
	}
	// check balance of supply
	if amount > pool.GetCash() {
		return fmt.Errorf("not enough token for withdraw. amount: %v, cash %v", amount, pool.GetCash())
	}

	// calcuate bill
	bill := amount * (pool.SupplyBill / pool.Supply)

	// update user account bill
	b.AccountDepositBills[user][symbol] -= bill

	// update pool
	pool.SupplyBill -= bill
	pool.Supply -= amount
	b.Pools[symbol] = pool

	return
}
