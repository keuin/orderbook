package orderbook

import (
	"fmt"
	"github.com/emirpasic/gods/trees/btree"
	"github.com/emirpasic/gods/utils"
	"slices"
)

func descendingIntComparator(a, b interface{}) int {
	aAsserted := a.(int)
	bAsserted := b.(int)
	switch {
	case aAsserted > bAsserted:
		return -1
	case aAsserted < bAsserted:
		return 1
	default:
		return 0
	}
}

func NewBTreeOrderBook() *BtreeOrderBook {
	return &BtreeOrderBook{
		bids: *btree.NewWith(8, descendingIntComparator),
		asks: *btree.NewWith(8, utils.IntComparator),
	}
}

type BtreeOrderBook struct {
	bids, asks btree.Tree
}

func (b *BtreeOrderBook) Snapshot(depth int) Snapshot {
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

func btreeUpdatePriceAmountAsc(tree *btree.Tree, price Price, amount int) {
	k := int(price)
	node := tree.GetNode(k)
	if node == nil {
		tree.Put(k, amount)
		return
	}
	i, ok := slices.BinarySearchFunc(node.Entries, k, func(entry *btree.Entry, price int) int {
		return entry.Key.(int) - price
	})
	if !ok {
		panic(fmt.Sprintf("node is found by btree container (%v), but is not present afterwards: %+v",
			k, node))
	}
	elem := node.Entries[i]
	amount += elem.Value.(int)
	if amount == 0 {
		tree.Remove(k)
		return
	}
	tree.Put(k, amount)
}

func btreeUpdatePriceAmountDesc(tree *btree.Tree, price Price, amount int) {
	k := int(price)
	node := tree.GetNode(k)
	if node == nil {
		tree.Put(k, amount)
		return
	}
	i, ok := slices.BinarySearchFunc(node.Entries, k, func(entry *btree.Entry, price int) int {
		return price - entry.Key.(int)
	})
	if !ok {
		panic(fmt.Sprintf("node is found by btree container (%v), but is not present afterwards: %+v",
			k, node))
	}
	elem := node.Entries[i]
	amount += elem.Value.(int)
	if amount == 0 {
		tree.Remove(k)
		return
	}
	tree.Put(k, amount)
}

func (b *BtreeOrderBook) Bid(price Price, amount int) {
	btreeUpdatePriceAmountDesc(&b.bids, price, amount)
}

func (b *BtreeOrderBook) Ask(price Price, amount int) {
	btreeUpdatePriceAmountAsc(&b.asks, price, amount)
}

func (b *BtreeOrderBook) GetBid(price Price) (amount int) {
	v, ok := b.bids.Get(int(price))
	if !ok {
		return 0
	}
	return v.(int)
}

func (b *BtreeOrderBook) GetAsk(price Price) (amount int) {
	v, ok := b.asks.Get(int(price))
	if !ok {
		return 0
	}
	return v.(int)
}
