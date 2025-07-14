package orderbook

import (
	"math/rand"
	"testing"
)

func BenchmarkArrayOrderBook(b *testing.B) {
	benchmarkOrderBook(b, &ArrayOrderBook{})
}

func BenchmarkNaiveOrderBook(b *testing.B) {
	benchmarkOrderBook(b, NewNaiveOrderBook())
}

func benchmarkOrderBook[T OrderBook](b *testing.B, ob T) {
	r := rand.NewSource(123456)
	for i := 0; i < b.N; i++ {
		isBid := r.Int63()&1 != 0
		isCancel := r.Int63()&1 != 0
		price := Price(uint32(r.Int63()) % 100)
		amount := int(r.Int63()) % 100
		if isCancel {
			amount = -amount
		}
		depth := int(r.Int63()%10) + 5
		_ = ob.Snapshot(depth)
		if isBid {
			if amount < 0 {
				currentAmount := ob.GetBid(price)
				if currentAmount == 0 {
					amount = -amount
				} else if amount+currentAmount < 0 {
					amount = -((-amount) % currentAmount)
				}
			}
			ob.Bid(price, amount)
		} else {
			if amount < 0 {
				currentAmount := ob.GetAsk(price)
				if currentAmount == 0 {
					amount = -amount
				} else if amount+currentAmount < 0 {
					amount %= currentAmount
				}
			}
			ob.Ask(price, amount)
		}
	}
}
