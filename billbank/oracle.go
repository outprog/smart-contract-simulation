package billbank

type Oracle struct {
	TokensPrice map[tSymbol]tPrice
}

func NewOracle() *Oracle {
	return &Oracle{map[tSymbol]tPrice{}}
}

func (o *Oracle) GetPrice(symbol string) float64 {
	if v, ok := o.TokensPrice[symbol]; ok {
		return v
	}
	return 0.0
}

func (o *Oracle) SetPrice(symbol string, price float64) {
	o.TokensPrice[symbol] = price
}
