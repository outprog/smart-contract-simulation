package mm

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	contract := NewMMContract(1000, 128000)
	daiCount := float64(0)
	for i := 0; i < 10; i++ {
		dai := contract.ETHtoDAI(10)
		daiCount += dai
		fmt.Println(contract)
		fmt.Println(dai)
	}
	fmt.Println(daiCount)

	ethCount := float64(0)
	for i := 0; i < 10; i++ {
		eth := contract.DAItoETH(daiCount / 10)
		ethCount += eth
		fmt.Println(contract)
		fmt.Println(eth)
	}
	fmt.Println(ethCount)

}
