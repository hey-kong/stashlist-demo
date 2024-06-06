package stashlist

import (
	"testing"

	"github.com/hey-kong/stashlist/cache/lru"
	"github.com/hey-kong/stashlist/cache/sieve"
	"github.com/hey-kong/stashlist/skiplist"
)

// LRU Put
func BenchmarkLruPutValue64B(b *testing.B) {
	lruCache = lru.New(b.N / 2)
	opLen := len(putOperations)

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		op := putOperations[n%opLen]
		lruCache.Add(op.key, op.value)
	}
}

// SIEVE Put
func BenchmarkSievePutValue64B(b *testing.B) {
	sieveCache = sieve.New(b.N / 2)
	opLen := len(putOperations)

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		op := putOperations[n%opLen]
		sieveCache.Add(op.key, op.value)
	}
}

// Skiplist Put
func BenchmarkSkiplistPutValue64B(b *testing.B) {
	l = skiplist.NewSkipList()
	opLen := len(putOperations)

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		op := putOperations[n%opLen]
		l.Add(op.key, op.value)
	}
}

func BenchmarkStashlistPutValue64B(b *testing.B) {
	myList = NewStashList()
	opLen := len(putOperations)

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		op := putOperations[n%opLen]
		myList.Add(op.key, op.value)
	}
}
