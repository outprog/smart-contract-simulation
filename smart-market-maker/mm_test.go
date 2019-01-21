package mm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestETHtoDAI(t *testing.T) {
	contract := NewMMContract(1000, 128000)

	ethAmount := 1.0
	daiAmount := 0.0

	// estimate
	estDAI := contract.EstimateETHtoDAI(ethAmount)

	// exchange
	daiAmount = contract.ETHtoDAI(ethAmount)
	assert.Equal(t, estDAI, daiAmount)

	assert.Equal(t, 1000+ethAmount, contract.BalanceOfETH)
	assert.Equal(t, 128000-daiAmount, contract.BalanceOfDAI)
	assert.True(t, true, contract.ProductValue > 1000*128000)
}

func TestDAItoETH(t *testing.T) {
	contract := NewMMContract(1000, 128000)

	ethAmount := 0.0
	daiAmount := 100.0

	// estimate
	estETH := contract.EstimateDAItoETH(daiAmount)

	// exchange
	ethAmount = contract.DAItoETH(daiAmount)
	assert.Equal(t, estETH, ethAmount)

	assert.Equal(t, 128000+daiAmount, contract.BalanceOfDAI)
	assert.Equal(t, 1000-ethAmount, contract.BalanceOfETH)
	assert.True(t, true, contract.ProductValue > 1000*128000)
}

func TestEstimate(t *testing.T) {
	contract := NewMMContract(1000, 128000)

	estETH1 := contract.EstimateDAItoETH(100)
	estETH2 := contract.EstimateDAItoETH(1000)
	assert.True(t, true, estETH1/100 > estETH2/1000)
	assert.Equal(t, float64(1000), contract.BalanceOfETH)
	assert.Equal(t, float64(128000), contract.BalanceOfDAI)

	estDAI1 := contract.EstimateETHtoDAI(1)
	estDAI2 := contract.EstimateETHtoDAI(2)
	assert.True(t, true, estDAI1/1 > estDAI2/2)
	assert.Equal(t, float64(1000), contract.BalanceOfETH)
	assert.Equal(t, float64(128000), contract.BalanceOfDAI)
}

func TestBigTrade(t *testing.T) {
	contract := NewMMContract(1000, 128000)

	contract.DAItoETH(1e30)
	assert.True(t, true, contract.BalanceOfETH > 0)
	contract.DAItoETH(1e90)
	assert.True(t, true, contract.BalanceOfETH > 0)
}

// func TestTrade(t *testing.T) {
// 	contract := NewMMContract(1000, 128000)
// 	daiCount := float64(0)
// 	for i := 0; i < 10; i++ {
// 		dai := contract.ETHtoDAI(10)
// 		daiCount += dai
// 		fmt.Println(contract)
// 		fmt.Println(dai)
// 	}
// 	fmt.Println(daiCount)

// 	ethCount := float64(0)
// 	for i := 0; i < 10; i++ {
// 		eth := contract.DAItoETH(daiCount / 10)
// 		ethCount += eth
// 		fmt.Println(contract)
// 		fmt.Println(eth)
// 	}
// 	fmt.Println(ethCount)
// }
