package orderbook

import (
	"github.com/emirpasic/gods/trees/btree"
	"github.com/emirpasic/gods/utils"
)

func NewBTreeRefOrderBook() *BtreeRefOrderBook {
	return &BtreeRefOrderBook{
		bids: btree.NewWith(8, descendingIntComparator),
		asks: btree.NewWith(8, utils.IntComparator),
	}
}

type BtreeRefOrderBook struct {
	bids, asks *btree.Tree
}

func (b BtreeRefOrderBook) Snapshot(depth int) Snapshot {
	var bestBids []PriceAmount
	if !b.bids.Empty() {
		bestBids = make([]PriceAmount, depth)
		it := b.bids.Iterator()
		it.Next()
		for i := 0; i < depth; i++ {
			bestBids[i] = PriceAmount{
				Price:  Price(it.Key().(int)),
				Amount: it.Value().(int),
			}
			if !it.Next() {
				bestBids = bestBids[:i+1]
				break
			}
		}
	}

	var bestAsks []PriceAmount
	if !b.asks.Empty() {
		bestAsks = make([]PriceAmount, depth)
		it := b.asks.Iterator()
		it.Next()
		for i := 0; i < depth; i++ {
			bestAsks[i] = PriceAmount{
				Price:  Price(it.Key().(int)),
				Amount: it.Value().(int),
			}
			if !it.Next() {
				bestAsks = bestAsks[:i+1]
				break
			}
		}
	}

	return Snapshot{
		BestBids: bestBids,
		BestAsks: bestAsks,
	}
}

func (b BtreeRefOrderBook) Bid(price Price, amount int) {
	btreeUpdatePriceAmountDesc(b.bids, price, amount)
}

func (b BtreeRefOrderBook) Ask(price Price, amount int) {
	btreeUpdatePriceAmountAsc(b.asks, price, amount)
}

func (b BtreeRefOrderBook) GetBid(price Price) (amount int) {
	v, ok := b.bids.Get(int(price))
	if !ok {
		return 0
	}
	return v.(int)
}

func (b BtreeRefOrderBook) GetAsk(price Price) (amount int) {
	v, ok := b.asks.Get(int(price))
	if !ok {
		return 0
	}
	return v.(int)
}
