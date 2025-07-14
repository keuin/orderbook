package orderbook

type OrderBook interface {
	Snapshot(depth int) Snapshot
	Bid(price Price, amount int)
	Ask(price Price, amount int)
}

type PriceAmount struct {
	Price  Price
	Amount int
}

type Snapshot struct {
	BestBids []PriceAmount
	BestAsks []PriceAmount
}

type Price int
