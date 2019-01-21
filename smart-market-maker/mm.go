package mm

import "fmt"

// MMContract is simulate of smart contract
// formula: (BalanceOfETH + ethAmount) * (BalanceOfDai - daiAmount) = BalanceOfETH * BalanceOfDAI
// @ProductValue = BalanceOfETH * BalanceOfDAI
// @increaseFactor
type MMContract struct {
	BalanceOfETH   float64
	BalanceOfDAI   float64
	ProductValue   float64
	increaseFactor float64
}

func NewMMContract(ethAmount, daiAmount float64) *MMContract {
	pV := ethAmount * daiAmount
	return &MMContract{
		BalanceOfETH:   ethAmount,
		BalanceOfDAI:   daiAmount,
		ProductValue:   pV,
		increaseFactor: 10,
	}
}

func (m *MMContract) String() string {
	return fmt.Sprintf("------------------------\n"+
		"BalanceOfETH: %v\n"+
		"BalanceOfDAI: %v\n"+
		"ProductValue: %v\n"+
		"CurrentValue: %v\n"+
		"Pirce:        %v\n"+
		"------------------------",
		m.BalanceOfETH,
		m.BalanceOfDAI,
		m.ProductValue,
		m.BalanceOfDAI+m.BalanceOfETH*m.EstimateETHtoDAI(1),
		m.EstimateETHtoDAI(1))
}

func (m *MMContract) ETHtoDAI(ethAmount float64) (daiAmount float64) {
	daiAmount = m.EstimateETHtoDAI(ethAmount)
	m.BalanceOfETH += ethAmount
	m.BalanceOfDAI -= daiAmount
	m.rebalance()
	return daiAmount
}
func (m *MMContract) DAItoETH(daiAmount float64) (ethAmount float64) {
	ethAmount = m.EstimateDAItoETH(daiAmount)
	m.BalanceOfDAI += daiAmount
	m.BalanceOfETH -= ethAmount
	m.rebalance()
	return ethAmount
}

func (m *MMContract) EstimateETHtoDAI(ethAmount float64) (daiAmount float64) {
	daiAmount = m.BalanceOfDAI - m.ProductValue/(m.BalanceOfETH+ethAmount)
	daiAmount -= daiAmount * m.increaseFactor / 100
	return daiAmount
}
func (m *MMContract) EstimateDAItoETH(daiAmount float64) (ethAmount float64) {
	ethAmount = m.BalanceOfETH - m.ProductValue/(m.BalanceOfDAI+daiAmount)
	ethAmount -= ethAmount * m.increaseFactor / 100
	return ethAmount
}

func (m *MMContract) rebalance() {
	newPv := m.BalanceOfETH * m.BalanceOfDAI
	if newPv > m.ProductValue {
		m.ProductValue = newPv
	}
}

func (m *MMContract) DepositETH(ethAmount float64) {
}
func (m *MMContract) WithdrawETH(ethAmount float64) {
}
func (m *MMContract) DepositDAI(daiAmount float64) {
}
func (m *MMContract) WithdrawDAI(daiAmount float64) {
}
