package money

import (
	"fmt"
)

type tUser = string
type tSymbol = string
type balance = float64

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
	SupplyBalances map[tUser]map[tSymbol]balance
	BorrowBalances map[tUser]map[tSymbol]balance

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
		SupplyBalances: map[tUser]map[tSymbol]balance{},
		BorrowBalances: map[tUser]map[tSymbol]balance{},
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
		if _, ok := supply[symbol]; ok {
			m.SupplyBalances[user][symbol] += amount
		} else {
			supply[symbol] = amount
		}
	} else {
		m.SupplyBalances[user] = map[tSymbol]float64{symbol: amount}
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
		if _, ok := borrow[symbol]; ok {
			m.BorrowBalances[user][symbol] += amount
		} else {
			borrow[symbol] = amount
		}
	} else {
		m.BorrowBalances[user] = map[tSymbol]float64{symbol: amount}
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
