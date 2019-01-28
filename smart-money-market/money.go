package money

import (
	"log"
)

type tUser = string
type tSymbol = string
type balance = float64

type Market struct {
	supply float64
	borrow float64
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

func (m *MoneyContract) Supply(amount float64, symbol string, user string) {
	if supply, ok := m.SupplyBalances[user]; ok {
		if _, ok := supply[symbol]; ok {
			m.SupplyBalances[user][symbol] += amount
		} else {
			supply[symbol] = amount
		}
	} else {
		m.SupplyBalances[user] = map[tSymbol]float64{symbol: amount}
	}

	m.Markets[symbol] = Market{
		supply: m.Markets[symbol].supply + amount,
		borrow: m.Markets[symbol].borrow,
	}
}

func (m *MoneyContract) Borrow(amount float64, symbol string, user string) {
	var market Market
	var ok bool
	if market, ok = m.Markets[symbol]; !ok {
		log.Printf("not support token: %v", symbol)
		return
	}

	if cash := market.supply - market.borrow; cash < amount {
		log.Printf("not enough cash: %v", cash)
		return
	}
}
