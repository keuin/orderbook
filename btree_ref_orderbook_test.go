package orderbook

import (
	"testing"
)

func FuzzBTreeRefOrderBook(f *testing.F) {
	fuzzOrderBook(f, NewBTreeRefOrderBook)
}
