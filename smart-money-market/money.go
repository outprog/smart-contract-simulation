package money

import (
	"fmt"
)

type tUser = string
type tSymbol = string
type balance = float64

// Market is token ledger
type Market struct {
	supply float64
	borrow float64
	growth float64
}

type MoneyContract struct {
	SupplyBalances map[tUser]map[tSymbol]balance
	BorrowBalances map[tUser]map[tSymbol]balance

	Markets        map[tSymbol]Market
	InitBorrowRate float64
}

func New() *MoneyContract {
	return &MoneyContract{
		SupplyBalances: map[tUser]map[tSymbol]balance{},
		BorrowBalances: map[tUser]map[tSymbol]balance{},
		Markets: map[tSymbol]Market{
			"ETH": Market{},
			"DAI": Market{},
		},
		InitBorrowRate: 0.05,
	}
}

func (m *MoneyContract) Supply(amount float64, symbol string, user string) error {
	// check market
	var market Market
	var ok bool
	if market, ok = m.Markets[symbol]; !ok {
		return fmt.Errorf("not support token: %v", symbol)
	}

	// update user supply
	if supply, ok := m.SupplyBalances[user]; ok {
		if _, ok := supply[symbol]; ok {
			m.SupplyBalances[user][symbol] += amount
		} else {
			supply[symbol] = amount
		}
	} else {
		m.SupplyBalances[user] = map[tSymbol]float64{symbol: amount}
	}

	// update market supply
	market.supply += amount
	m.Markets[symbol] = market

	return nil
}

func (m *MoneyContract) Borrow(amount float64, symbol string, user string) error {
	// check market
	var market Market
	var ok bool
	if market, ok = m.Markets[symbol]; !ok {
		return fmt.Errorf("not support token: %v", symbol)
	}

	// check cash
	if cash := market.supply - market.borrow; cash < amount {
		return fmt.Errorf("not enough cash: %v", cash)
	}

	// update user borrow
	if borrow, ok := m.BorrowBalances[user]; ok {
		if _, ok := borrow[symbol]; ok {
			m.BorrowBalances[user][symbol] += amount
		} else {
			borrow[symbol] = amount
		}
	} else {
		m.BorrowBalances[user] = map[tSymbol]float64{symbol: amount}
	}

	// update market borrow
	market.borrow += amount
	m.Markets[symbol] = market

	return nil
}
