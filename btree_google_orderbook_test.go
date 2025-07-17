package orderbook

import (
	"testing"
)

func FuzzBTreeGoogleOrderBook(f *testing.F) {
	fuzzOrderBook(f, NewBTreeGoogleOrderBook)
}
