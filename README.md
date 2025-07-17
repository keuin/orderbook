# Orderbook

Limit Order Book implementations backed by different data structures:

- Naive HashMap
- Sorted array
- BTree (non-generic and generic versions from `github.com/emirpasic/gods`, `github.com/google/btree`)

## Interface

The `OrderBook` interface is defined in `orderbook.go`:

```go
package orderbook


type OrderBook interface {
	Snapshot(depth int) Snapshot
	Bid(price Price, amount int)
	Ask(price Price, amount int)
	GetBid(price Price) (amount int)
	GetAsk(price Price) (amount int)
}

type Price int

type PriceAmount struct {
	Price  Price
	Amount int
}

type Snapshot struct {
	BestBids []PriceAmount
	BestAsks []PriceAmount
}
```

## Benchmark

```
goos: linux
goarch: amd64
pkg: github.com/keuin/orderbook
cpu: 12th Gen Intel(R) Core(TM) i9-12900K
BenchmarkArrayOrderBook-24               8621535               133.6 ns/op
BenchmarkNaiveOrderBook-24                146120              8072 ns/op
BenchmarkBTreeOrderBook-24               2425657               473.4 ns/op
BenchmarkBTreeRefOrderBook-24            2425458               495.5 ns/op
BenchmarkBTreeGoogleOrderBook-24         4143693               286.9 ns/op
BenchmarkBTreeGoogleNoGOrderBook-24      2424078               493.5 ns/op
PASS
ok      github.com/keuin/orderbook      9.973s
```

## TODO

- More benchmark scenarios (e.g. drastically price shifting)
- Reduce redundant search operations in BTree using custom implementation
