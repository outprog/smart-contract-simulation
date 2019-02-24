package billbank

import (
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
