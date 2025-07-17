package orderbook

import (
	"testing"
)

func FuzzBTreeGoogleNoGOrderBook(f *testing.F) {
	fuzzOrderBook(f, NewBTreeGoogleNoGOrderBook)
}
