package billbank

import (
	"log"
	"math"
)

type tUser = string
type tSymbol = string
type tBill = float64
type tPrice = float64

type TokenPool struct {
	SupplyBill float64
	Supply     float64
	BorrowBill float64
	Borrow     float64

	// last liquidate blockNumber
	liquidateIndex uint64
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

	Oralcer *Oracle

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
		Oralcer:     NewOracle(),
		BlockNumber: 1,
		borrowRate:  0.01,
	}
}

func (b *Billbank) liquidate(symbol string) {
	pool := b.getPool(symbol)

	growth := b.calculateGrowth(symbol)
	pool.Supply += growth
	pool.Borrow += growth

	// update pool
	pool.liquidateIndex = b.BlockNumber
	b.Pools[symbol] = pool
}

func (b *Billbank) calculateGrowth(symbol string) float64 {
	pool := b.getPool(symbol)

	growth := 0.0
	borrow := pool.Borrow
	if borrow != 0.0 {
		// Compound interest
		// formula:
		//		b: borrow
		//		r: rate
		//		n: block number
		//		b = b * (1+r)^n
		borrow = borrow * math.Pow(
			1.0+b.borrowRate,
			float64(b.BlockNumber-pool.liquidateIndex),
		)
		growth = borrow - pool.Borrow
	}
	return growth
}

func (b *Billbank) getPool(symbol string) (pool TokenPool) {
	var ok bool
	if pool, ok = b.Pools[symbol]; !ok {
		log.Panicf("not support token: %v", symbol)
	}
	return
}

func (b *Billbank) NetValueOf(user string) float64 {
	supplyValue := 0.0
	if acc, ok := b.AccountDepositBills[user]; ok {
		for sym, bill := range acc {
			if bill != 0.0 {
				supplyValue += b.SupplyValueOf(sym, user)
			}
		}
	}

	borrowValue := 0.0
	if acc, ok := b.AccountBorrowBills[user]; ok {
		for sym, bill := range acc {
			if bill != 0.0 {
				borrowValue += b.BorrowValueOf(sym, user)
			}
		}
	}

	return supplyValue - borrowValue
}
