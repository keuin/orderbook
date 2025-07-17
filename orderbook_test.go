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

func BenchmarkBTreeOrderBook(b *testing.B) {
	benchmarkOrderBook(b, NewBTreeOrderBook())
}

func BenchmarkBTreeRefOrderBook(b *testing.B) {
	benchmarkOrderBook(b, NewBTreeRefOrderBook())
}

func BenchmarkBTreeGoogleOrderBook(b *testing.B) {
	benchmarkOrderBook(b, NewBTreeGoogleOrderBook())
}

func BenchmarkBTreeGoogleNoGOrderBook(b *testing.B) {
	benchmarkOrderBook(b, NewBTreeGoogleNoGOrderBook())
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
			if amount == 0 {
				continue
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
			if amount == 0 {
				continue
			}
			ob.Ask(price, amount)
		}
	}
}

func fuzzOrderBook[T OrderBook](f *testing.F, newOrderBook func() T) {
	f.Fuzz(func(t *testing.T, seed int64) {
		r := rand.NewSource(seed)
		refOb := NewNaiveOrderBook()
		ob := newOrderBook()
		for i := 0; i < 10000; i++ {
			isBid := r.Int63()&1 != 0
			isCancel := r.Int63()&1 != 0
			price := Price(uint32(r.Int63()) % 100)
			amount := int(r.Int63()) % 100
			if isCancel {
				amount = -amount
			}
			depth := int(r.Int63()%10) + 5
			snapshot := ob.Snapshot(depth)
			refSnapshot := refOb.Snapshot(depth)
			if !snapshot.Equals(refSnapshot) {
				t.Fatal("incorrect snapshot, expected: ", refSnapshot, " got: ", snapshot)
			}
			if isBid {
				if amount < 0 {
					currentAmount := ob.GetBid(price)
					if currentAmount == 0 {
						amount = -amount
					} else if amount+currentAmount < 0 {
						amount = -((-amount) % currentAmount)
					}
				}
				if amount == 0 {
					continue
				}
				t.Log("BID", price, amount)
				ob.Bid(price, amount)
				refOb.Bid(price, amount)
			} else {
				if amount < 0 {
					currentAmount := ob.GetAsk(price)
					if currentAmount == 0 {
						amount = -amount
					} else if amount+currentAmount < 0 {
						amount %= currentAmount
					}
				}
				if amount == 0 {
					continue
				}
				t.Log("ASK", price, amount)
				ob.Ask(price, amount)
				refOb.Ask(price, amount)
			}
		}
	})
}
