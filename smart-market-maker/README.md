# smart market maker

Base on formula: x*y=k

```
(BalanceOfETH + ethAmount) * (BalanceOfDai - daiAmount) = BalanceOfETH * BalanceOfDAI
```

Added `PositionOffset` initialize a price in contract.

## REFERENCE

- [x*y=k market makers](https://ethresear.ch/t/improving-front-running-resistance-of-x-y-k-market-makers/1281)
- [Uniswap Whitepaper](https://hackmd.io/C-DvwDSfSxuh-Gd4WKE_ig)
- [Uniswap — A Unique Exchange](https://medium.com/@cyrus.younessi/uniswap-a-unique-exchange-f4ef44f807bf)
- [没有市场深度，即买即卖的去中心化交易所](https://www.jianshu.com/p/9a86a9252f9b)
