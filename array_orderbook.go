package orderbook

import (
	"slices"
)

type ArrayOrderBook struct {
	bids, asks binSearchBestPrices
}

func (a *ArrayOrderBook) Snapshot(depth int) Snapshot {
	ret := Snapshot{
		BestBids: a.bids,
		BestAsks: a.asks,
	}
	if len(ret.BestBids) > depth {
		ret.BestBids = ret.BestBids[:depth]
	}
	if len(ret.BestAsks) > depth {
		ret.BestAsks = ret.BestAsks[:depth]
	}
	ret = Snapshot{
		BestBids: slices.Clone(ret.BestBids),
		BestAsks: slices.Clone(ret.BestAsks),
	}
	return ret
}

func (a *ArrayOrderBook) Bid(price Price, amount int) {
	a.bids.UpdateBestPriceAmounts(price, amount, false)
}

func (a *ArrayOrderBook) Ask(price Price, amount int) {
	a.asks.UpdateBestPriceAmounts(price, amount, true)
}

func (a *ArrayOrderBook) GetBid(price Price) (amount int) {
	return a.bids.GetAmount(price)
}

func (a *ArrayOrderBook) GetAsk(price Price) (amount int) {
	return a.asks.GetAmount(price)
}

type binSearchBestPrices []PriceAmount

func (p *binSearchBestPrices) GetAmount(price Price) (amount int) {
	i, ok := slices.BinarySearchFunc(*p, price, func(pa PriceAmount, price Price) int {
		return int(pa.Price - price)
	})
	if !ok {
		return 0
	}
	return (*p)[i].Amount
}

func (p *binSearchBestPrices) UpdateBestPriceAmounts(price Price, amount int, ask bool) {
	const maxLinearSearchCount = 5
	arr := *p
	n := min(maxLinearSearchCount, len(arr))
	for i := 0; i < n; i++ {
		if arr[i].Price != price {
			continue
		}
		amount += arr[i].Amount
		if amount < 0 {
			panic("invalid order: negative resulting amount")
		}
		if amount == 0 {
			*p = slices.Delete(arr, i, i+1)
		} else {
			arr[i].Amount = amount
		}
		return
	}
	i, j := 0, len(arr)-1
	for i <= j {
		mid := i + (j-i)/2
		midVal := arr[mid].Price
		if price == midVal {
			// price found, update in-place
			amount += arr[mid].Amount
			if amount < 0 {
				panic("invalid order: negative resulting amount")
			}
			if amount == 0 {
				// no orders in this price, remove price item
				*p = slices.Delete(arr, mid, mid+1)
				return
			}
			arr[mid].Amount = amount
			return
		} else if price < midVal == ask {
			// proceed on the left part
			j = mid - 1
		} else if price > midVal == ask {
			// proceed on the right part
			i = mid + 1
		}
	}
	// not found, insert
	if amount == 0 {
		return
	}
	if len(arr) == 0 {
		*p = make(binSearchBestPrices, 0, 128)
		*p = append(*p, PriceAmount{Price: price, Amount: amount})
		return
	}
	arr = append(arr, PriceAmount{})
	copy(arr[i+1:], arr[i:])
	arr[i] = PriceAmount{Price: price, Amount: amount}
	*p = arr
}
