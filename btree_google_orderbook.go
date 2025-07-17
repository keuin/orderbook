package orderbook

import (
	"github.com/google/btree"
)

func NewBTreeGoogleOrderBook() *BtreeGoogleOrderBook {
	return &BtreeGoogleOrderBook{
		bids: *btree.NewG(8, func(a, b PriceAmount) bool {
			return a.Price > b.Price
		}),
		asks: *btree.NewG(8, func(a, b PriceAmount) bool {
			return a.Price < b.Price
		}),
	}
}

type BtreeGoogleOrderBook struct {
	bids, asks btree.BTreeG[PriceAmount]
}

func (b *BtreeGoogleOrderBook) Snapshot(depth int) Snapshot {
	var bestBids []PriceAmount
	if b.bids.Len() != 0 {
		bestBids = make([]PriceAmount, depth)
		i := 0
		b.bids.Ascend(func(item PriceAmount) bool {
			bestBids[i] = item
			i++
			return i != depth
		})
		bestBids = bestBids[:i]
	}

	var bestAsks []PriceAmount
	if b.asks.Len() != 0 {
		bestAsks = make([]PriceAmount, depth)
		i := 0
		b.asks.Ascend(func(item PriceAmount) bool {
			bestAsks[i] = item
			i++
			return i != depth
		})
		bestAsks = bestAsks[:i]
	}

	return Snapshot{
		BestBids: bestBids,
		BestAsks: bestAsks,
	}
}

func btreeGoogleUpdatePriceAmount(tree *btree.BTreeG[PriceAmount], price Price, amount int) {
	newItem, ok := tree.Get(PriceAmount{Price: price})
	if ok {
		newItem.Amount += amount
		if newItem.Amount == 0 {
			tree.Delete(PriceAmount{Price: price})
			return
		} else if newItem.Amount < 0 {
			panic("invalid amount")
		}
	} else {
		newItem = PriceAmount{Price: price, Amount: amount}
	}
	tree.ReplaceOrInsert(newItem)
}

func (b *BtreeGoogleOrderBook) Bid(price Price, amount int) {
	btreeGoogleUpdatePriceAmount(&b.bids, price, amount)
}

func (b *BtreeGoogleOrderBook) Ask(price Price, amount int) {
	btreeGoogleUpdatePriceAmount(&b.asks, price, amount)
}

func (b *BtreeGoogleOrderBook) GetBid(price Price) (amount int) {
	v, ok := b.bids.Get(PriceAmount{Price: price})
	if !ok {
		return 0
	}
	return v.Amount
}

func (b *BtreeGoogleOrderBook) GetAsk(price Price) (amount int) {
	v, ok := b.asks.Get(PriceAmount{Price: price})
	if !ok {
		return 0
	}
	return v.Amount
}
