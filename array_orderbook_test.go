package orderbook

import (
	"testing"
)

func FuzzArrayOrderBook(f *testing.F) {
	fuzzOrderBook(f, func() *ArrayOrderBook {
		return &ArrayOrderBook{}
	})
}
