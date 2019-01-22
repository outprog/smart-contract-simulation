package mm

import (
	"fmt"
	"math"
)

// MMContract is simulate of smart contract
// formula: x*y=k (BalanceOfETH + ethAmount) * (BalanceOfDai - daiAmount) = BalanceOfETH * BalanceOfDAI
// @ProductValue = BalanceOfETH * BalanceOfDAI
// @growthFactor
// @positionOffset
type MMContract struct {
	BalanceOfETH   float64
	BalanceOfDAI   float64
	ProductValue   float64
	growthFactor   float64
	positionOffset float64
}

// NewMMContract new a contract
// @ethAmount
// @daiAmount
// @initPrice is init price, price is daiAmount/ethAmount
func NewMMContract(ethAmount, daiAmount, initPrice float64) *MMContract {
	product := ethAmount * daiAmount
	offset := float64(1)
	if initPrice > 0 {
		// x*y=k    BalanceOfETH * BalanceOfDAI = ProductValue
		// y/x=p    BalanceOfDAI / BalanceOfETH = Price
		// derivation:
		// ix*iy=k
		// => iy=k/ix
		// iy/ix=ip (iBalanceOfDAI / iBalanceOfETH = initPrice)
		// => (k/ix)/ix=ip
		// => ix=(k/ip)^(1/2)
		// offset = ix/x
		// => offset = (k/ip)^(1/2)/x
		offset = math.Sqrt(product/initPrice) / ethAmount
	}
	return &MMContract{
		BalanceOfETH:   ethAmount,
		BalanceOfDAI:   daiAmount,
		ProductValue:   product,
		growthFactor:   3,
		positionOffset: offset,
	}
}

func (m *MMContract) String() string {
	return fmt.Sprintf("---------\n"+
		"BalanceOfETH:         %v\n"+
		"BalanceOfDAI:         %v\n"+
		"ProductValue:         %v\n"+
		"PositionOffset:       %v\n"+
		"CurrentPirce:         %v\n"+
		"CurrentValue(DAI):    %v\n"+
		"CurrentValue(ETH):    %v\n"+
		"---------------------------------------------",
		m.BalanceOfETH,
		m.BalanceOfDAI,
		m.ProductValue,
		m.positionOffset,
		m.Price(),
		m.BalanceOfDAI+m.BalanceOfETH*m.Price(),
		m.BalanceOfETH+m.BalanceOfDAI/m.Price(),
	)
}

func (m *MMContract) Price() float64 {
	eth, dai := m.balanceWithOffset()
	return dai / eth
}

func (m *MMContract) ETHtoDAI(ethAmount float64) (daiAmount float64) {
	daiAmount = m.EstimateETHtoDAI(ethAmount)

	// Insufficient
	if m.BalanceOfDAI < daiAmount {
		return -1
	}

	m.BalanceOfETH += ethAmount
	m.BalanceOfDAI -= daiAmount
	m.growth()
	return daiAmount
}
func (m *MMContract) DAItoETH(daiAmount float64) (ethAmount float64) {
	ethAmount = m.EstimateDAItoETH(daiAmount)

	// Insufficient
	if m.BalanceOfETH < ethAmount {
		return -1
	}

	m.BalanceOfDAI += daiAmount
	m.BalanceOfETH -= ethAmount
	m.growth()
	return ethAmount
}

func (m *MMContract) EstimateETHtoDAI(ethAmount float64) (daiAmount float64) {
	eth, dai := m.balanceWithOffset()
	daiAmount = dai - m.ProductValue/(eth+ethAmount)
	daiAmount -= daiAmount * m.growthFactor / 1000
	return daiAmount
}
func (m *MMContract) EstimateDAItoETH(daiAmount float64) (ethAmount float64) {
	eth, dai := m.balanceWithOffset()
	ethAmount = eth - m.ProductValue/(dai+daiAmount)
	ethAmount -= ethAmount * m.growthFactor / 1000
	return ethAmount
}

func (m *MMContract) balanceWithOffset() (ethAmount, daiAmount float64) {
	ethAmount = m.BalanceOfETH * m.positionOffset
	daiAmount = m.ProductValue / ethAmount
	return
}

func (m *MMContract) growth() {
	newPv := m.BalanceOfETH * m.BalanceOfDAI
	if newPv > m.ProductValue {
		m.ProductValue = newPv
	}
}

func (m *MMContract) resetOffset(price float64) {
	m.positionOffset = math.Sqrt(m.ProductValue/price) / m.BalanceOfETH
}

func (m *MMContract) DepositETH(ethAmount float64) {
	price := m.Price()
	m.BalanceOfETH += ethAmount
	m.growth()
	m.resetOffset(price)
}
func (m *MMContract) WithdrawETH(ethAmount float64) {
	if ethAmount > m.BalanceOfETH {
		return
	}
	price := m.Price()
	m.BalanceOfETH -= ethAmount
	m.ProductValue = m.BalanceOfDAI * m.BalanceOfETH
	m.resetOffset(price)
}
func (m *MMContract) DepositDAI(daiAmount float64) {
	price := m.Price()
	m.BalanceOfDAI += daiAmount
	m.growth()
	m.resetOffset(price)
}
func (m *MMContract) WithdrawDAI(daiAmount float64) {
	if daiAmount > m.BalanceOfDAI {
		return
	}
	price := m.Price()
	m.BalanceOfDAI -= daiAmount
	m.ProductValue = m.BalanceOfDAI * m.BalanceOfETH
	m.resetOffset(price)
}
