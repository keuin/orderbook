package orderbook

import (
	"testing"
)

func FuzzBTreeOrderBook(f *testing.F) {
	fuzzOrderBook(f, NewBTreeOrderBook)
}
