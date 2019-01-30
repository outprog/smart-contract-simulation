package money

import (
	"fmt"
)

type tUser = string
type tSymbol = string

// Balance is recent balance after most recent balance-changing action
// @Principal total balance with accrued interest after applying the customer's most recent balance-changing action
type Balance struct {
	Principal   float64
	BlockNumber uint64
}

// RateCollection saved rate event with block number
type RateCollection struct {
	SupplyRate  float64
	BorrowRate  float64
	BlockNumber uint64
}

// Market is token ledger
type Market struct {
	Supply float64
	Borrow float64

	BaseRate float64
	RateList []RateCollection
}

// MoneyContact
// @BlockNumber is simulate of current block number
type MoneyContract struct {
	SupplyBalances map[tUser]map[tSymbol]Balance
	BorrowBalances map[tUser]map[tSymbol]Balance

	Markets map[tSymbol]Market

	BlockNumber uint64
}

func New() *MoneyContract {
	ethMarket := Market{
		RateList: []RateCollection{},
		BaseRate: 0.05,
	}
	daiMarket := Market{
		RateList: []RateCollection{},
		BaseRate: 0.05,
	}
	return &MoneyContract{
		SupplyBalances: map[tUser]map[tSymbol]Balance{},
		BorrowBalances: map[tUser]map[tSymbol]Balance{},
		Markets: map[tSymbol]Market{
			"ETH": ethMarket,
			"DAI": daiMarket,
		},
		BlockNumber: 1,
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
		if balance, ok := supply[symbol]; ok {
			balance.Principal += amount
			m.SupplyBalances[user][symbol] = balance
		} else {
			supply[symbol] = Balance{amount, m.BlockNumber}
		}
	} else {
		m.SupplyBalances[user] = map[tSymbol]Balance{symbol: Balance{amount, m.BlockNumber}}
	}

	// update market supply
	market.Supply += amount
	m.Markets[symbol] = market

	return m.calculateRate(symbol)
}

func (m *MoneyContract) Borrow(amount float64, symbol string, user string) error {
	// check market
	var market Market
	var ok bool
	if market, ok = m.Markets[symbol]; !ok {
		return fmt.Errorf("not support token: %v", symbol)
	}

	// check cash
	if cash := market.Supply - market.Borrow; cash < amount {
		return fmt.Errorf("not enough cash: %v", cash)
	}

	// update user borrow
	if borrow, ok := m.BorrowBalances[user]; ok {
		if balance, ok := borrow[symbol]; ok {
			balance.Principal += amount
			m.BorrowBalances[user][symbol] = balance
		} else {
			borrow[symbol] = Balance{amount, m.BlockNumber}
		}
	} else {
		m.BorrowBalances[user] = map[tSymbol]Balance{symbol: Balance{amount, m.BlockNumber}}
	}

	// update market borrow
	market.Borrow += amount
	m.Markets[symbol] = market

	return m.calculateRate(symbol)
}

func (m *MoneyContract) calculateRate(symbol string) error {
	// check market
	var market Market
	var ok bool
	if market, ok = m.Markets[symbol]; !ok {
		return fmt.Errorf("not support token: %v", symbol)
	}

	// init rate collection
	rateCollection := RateCollection{
		SupplyRate:  0,
		BorrowRate:  0,
		BlockNumber: m.BlockNumber,
	}

	// had borrow, calculate rate
	if market.Borrow != 0.0 {
		rateCollection.BorrowRate = market.BaseRate
		rateCollection.SupplyRate = (market.Borrow * rateCollection.BorrowRate) / market.Supply
	}

	market.RateList = append(market.RateList, rateCollection)
	m.Markets[symbol] = market
	return nil
}
