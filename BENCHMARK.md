# benchmark

## 100 prices, 8-order btree

```
$ go test -bench=. -benchtime=20s
goos: linux
goarch: amd64
pkg: github.com/keuin/orderbook
cpu: 12th Gen Intel(R) Core(TM) i9-12900K
BenchmarkArrayOrderBook-24      181495662              132.4 ns/op
BenchmarkNaiveOrderBook-24       2969110              8167 ns/op
BenchmarkBTreeOrderBook-24      49960929               480.3 ns/op
PASS
ok      github.com/keuin/orderbook      95.060s
```

## 10000 prices, 8-order btree

```
$ go test -bench=. -benchtime=20s # 10000 unique prices
goos: linux
goarch: amd64
pkg: github.com/keuin/orderbook
cpu: 12th Gen Intel(R) Core(TM) i9-12900K
BenchmarkArrayOrderBook-24      120053476              199.0 ns/op
BenchmarkNaiveOrderBook-24         88822           1085979 ns/op
BenchmarkBTreeOrderBook-24      25785922               944.6 ns/op
PASS
ok      github.com/keuin/orderbook      169.347s
```

## 10000 prices, 128-order btree

```
$ go test -bench=. -benchtime=20s # 10000 unique prices, 128-order btree
goos: linux
goarch: amd64
pkg: github.com/keuin/orderbook
cpu: 12th Gen Intel(R) Core(TM) i9-12900K
BenchmarkArrayOrderBook-24      100000000              200.3 ns/op
BenchmarkNaiveOrderBook-24         89763           1090439 ns/op
BenchmarkBTreeOrderBook-24      24580310               995.6 ns/op
PASS
ok      github.com/keuin/orderbook      147.025s
```
