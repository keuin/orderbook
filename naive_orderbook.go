package orderbook

import "slices"

func NewNaiveOrderBook() *NaiveOrderBook {
	return &NaiveOrderBook{
		bids: make(map[Price]int, 100),
		asks: make(map[Price]int, 100),
	}
}

type NaiveOrderBook struct {
	bids, asks map[Price]int
}

func (b *NaiveOrderBook) Snapshot(depth int) Snapshot {
	asks := make([]PriceAmount, len(b.asks))
	var i int
	for price, amount := range b.asks {
		asks[i] = PriceAmount{
			Price:  price,
			Amount: amount,
		}
		i++
	}
	slices.SortFunc(asks, func(a, b PriceAmount) int {
		return int(a.Price - b.Price)
	})
	if depth < len(asks) {
		asks = asks[:depth]
	}

	i = 0
	bids := make([]PriceAmount, len(b.bids))
	for price, amount := range b.bids {
		bids[i] = PriceAmount{
			Price:  price,
			Amount: amount,
		}
		i++
	}
	slices.SortFunc(bids, func(a, b PriceAmount) int {
		return int(b.Price - a.Price)
	})
	if depth < len(bids) {
		bids = bids[:depth]
	}
	return Snapshot{
		BestBids: bids,
		BestAsks: asks,
	}
}

func (b *NaiveOrderBook) Bid(price Price, amount int) {
	b.bids[price] += amount
	if b.bids[price] < 0 {
		panic("invalid argument: negative amount")
	}
	if b.bids[price] == 0 {
		delete(b.bids, price)
	}
}

func (b *NaiveOrderBook) Ask(price Price, amount int) {
	b.asks[price] += amount
	if b.asks[price] < 0 {
		panic("invalid argument: negative amount")
	}
	if b.asks[price] == 0 {
		delete(b.asks, price)
	}
}

func (b *NaiveOrderBook) GetBid(price Price) (amount int) {
	return b.bids[price]
}

func (b *NaiveOrderBook) GetAsk(price Price) (amount int) {
	return b.asks[price]
}
