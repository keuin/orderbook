package orderbook

import (
	"github.com/google/btree"
)

type btreeGoogleNoGItemAsc PriceAmount

func (b btreeGoogleNoGItemAsc) Less(than btree.Item) bool {
	return b.Price < than.(btreeGoogleNoGItemAsc).Price
}

type btreeGoogleNoGItemDesc PriceAmount

func (b btreeGoogleNoGItemDesc) Less(than btree.Item) bool {
	return b.Price > than.(btreeGoogleNoGItemDesc).Price
}

func NewBTreeGoogleNoGOrderBook() *BtreeGoogleNoGOrderBook {
	return &BtreeGoogleNoGOrderBook{
		bids: *btree.New(8),
		asks: *btree.New(8),
	}
}

type BtreeGoogleNoGOrderBook struct {
	bids, asks btree.BTree
}

func (b *BtreeGoogleNoGOrderBook) Snapshot(depth int) Snapshot {
	var bestBids []PriceAmount
	if b.bids.Len() != 0 {
		bestBids = make([]PriceAmount, depth)
		i := 0
		b.bids.Ascend(func(item btree.Item) bool {
			bestBids[i] = PriceAmount(item.(btreeGoogleNoGItemDesc))
			i++
			return i != depth
		})
		bestBids = bestBids[:i]
	}

	var bestAsks []PriceAmount
	if b.asks.Len() != 0 {
		bestAsks = make([]PriceAmount, depth)
		i := 0
		b.asks.Ascend(func(item btree.Item) bool {
			bestAsks[i] = PriceAmount(item.(btreeGoogleNoGItemAsc))
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

func (b *BtreeGoogleNoGOrderBook) Bid(price Price, amount int) {
	newItem := b.bids.Get(btreeGoogleNoGItemDesc{Price: price})
	if newItem != nil {
		newItem := newItem.(btreeGoogleNoGItemDesc)
		newItem.Amount += amount
		if newItem.Amount == 0 {
			b.bids.Delete(btreeGoogleNoGItemDesc{Price: price})
			return
		} else if newItem.Amount < 0 {
			panic("invalid amount")
		}
		b.bids.ReplaceOrInsert(newItem)
	} else {
		b.bids.ReplaceOrInsert(btreeGoogleNoGItemDesc{Price: price, Amount: amount})
	}
}

func (b *BtreeGoogleNoGOrderBook) Ask(price Price, amount int) {
	newItem := b.asks.Get(btreeGoogleNoGItemAsc{Price: price})
	if newItem != nil {
		newItem := newItem.(btreeGoogleNoGItemAsc)
		newItem.Amount += amount
		if newItem.Amount == 0 {
			b.asks.Delete(btreeGoogleNoGItemAsc{Price: price})
			return
		} else if newItem.Amount < 0 {
			panic("invalid amount")
		}
		b.asks.ReplaceOrInsert(newItem)
	} else {
		b.asks.ReplaceOrInsert(btreeGoogleNoGItemAsc{Price: price, Amount: amount})
	}
}

func (b *BtreeGoogleNoGOrderBook) GetBid(price Price) (amount int) {
	v := b.bids.Get(btreeGoogleNoGItemDesc{Price: price})
	if v == nil {
		return 0
	}
	return v.(btreeGoogleNoGItemDesc).Amount
}

func (b *BtreeGoogleNoGOrderBook) GetAsk(price Price) (amount int) {
	v := b.asks.Get(btreeGoogleNoGItemAsc{Price: price})
	if v == nil {
		return 0
	}
	return v.(btreeGoogleNoGItemAsc).Amount
}
