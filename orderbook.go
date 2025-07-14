package orderbook

type OrderBook interface {
	Snapshot(depth int) Snapshot
	Bid(price Price, amount int)
	Ask(price Price, amount int)
	GetBid(price Price) (amount int)
	GetAsk(price Price) (amount int)
}

type PriceAmount struct {
	Price  Price
	Amount int
}

type Snapshot struct {
	BestBids []PriceAmount
	BestAsks []PriceAmount
}

func (s Snapshot) Equals(other Snapshot) bool {
	if len(s.BestBids) != len(other.BestBids) {
		return false
	}
	if len(s.BestAsks) != len(other.BestAsks) {
		return false
	}
	for i := range s.BestBids {
		if s.BestBids[i] != other.BestBids[i] {
			return false
		}
	}
	for i := range s.BestAsks {
		if s.BestAsks[i] != other.BestAsks[i] {
			return false
		}
	}
	return true
}

type Price int
