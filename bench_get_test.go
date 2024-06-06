package stashlist

import (
	"testing"
)

// LRU Get
func BenchmarkLruGetValue64B(b *testing.B) {
	initLruCache(b.N)
	opLen := len(getOperations)

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		op := getOperations[n%opLen]
		lruCache.Get(op.key)
	}
}

// SIEVE Get
func BenchmarkSieveGetValue64B(b *testing.B) {
	initSieveCache(b.N)
	opLen := len(getOperations)

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		op := getOperations[n%opLen]
		sieveCache.Get(op.key)
	}
}

// Skiplist Get
func BenchmarkSkiplistGetValue64B(b *testing.B) {
	initSkiplist(b.N)
	opLen := len(getOperations)

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		op := getOperations[n%opLen]
		l.Get(op.key)
	}
}

func BenchmarkStashlistGetValue64B(b *testing.B) {
	initStashlist(b.N)
	opLen := len(getOperations)

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		op := getOperations[n%opLen]
		myList.Get(op.key)
	}
}
