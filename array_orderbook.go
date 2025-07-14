package orderbook

import (
	"fmt"
	"slices"
)

type ArrayOrderBook struct {
	bids, asks bestPrices
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
	a.bids.updateBestPriceAmounts(price, amount)
}

func (a *ArrayOrderBook) Ask(price Price, amount int) {
	a.asks.updateBestPriceAmounts(price, amount)
}

type bestPrices []PriceAmount

func (p *bestPrices) updateBestPriceAmounts(price Price, amount int) {
	arr := *p
	i, j := 0, len(arr)-1
	for i <= j {
		mid := i + (j-i)/2
		midVal := arr[mid].Price
		if price < midVal {
			// proceed on the left part
			j = mid - 1
		} else if price > midVal {
			// proceed on the right part
			i = mid + 1
		} else {
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
		}
		if arr[i].Price < price && arr[j].Price > price {
			if j-i != 1 {
				panic(fmt.Sprintf("invalid state, i=%v, j=%v", i, j))
			}
			*p = slices.Insert(arr, i+1, PriceAmount{
				Price:  price,
				Amount: amount,
			})
			return
		}
	}
}
