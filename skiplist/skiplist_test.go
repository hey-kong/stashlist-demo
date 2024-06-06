package skiplist

import (
	"bytes"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"unsafe"
)

var benchList *SkipList

func IntToBytes(n int) []byte {
	return []byte(strconv.Itoa(n))
}

func init() {
	// Initialize a big SkipList for the Get() benchmark
	benchList = NewSkipList()

	for i := 0; i <= 10000000; i++ {
		benchList.Add(strconv.Itoa(i), []byte{})
	}

	// Display the sizes of our basic structs
	var sl SkipList
	var el Element
	fmt.Printf("Structure sizes: SkipList is %v, Element is %v bytes\n", unsafe.Sizeof(sl), unsafe.Sizeof(el))
}

func checkSanity(list *SkipList, t *testing.T) {
	// each level must be correctly ordered
	for k, v := range list.next {
		//t.Log("Level", k)

		if v == nil {
			continue
		}

		if k > len(v.next) {
			t.Fatal("first node's level must be no less than current level")
		}

		next := v
		cnt := 1

		for next.next[k] != nil {
			if !(next.next[k].key >= next.key) {
				t.Fatalf("next key value must be greater than prev key value. [next:%v] [prev:%v]", next.next[k].key, next.key)
			}

			if k > len(next.next) {
				t.Fatalf("node's level must be no less than current level. [cur:%v] [node:%v]", k, next.next)
			}

			next = next.next[k]
			cnt++
		}

		if k == 0 {
			if cnt != list.Length {
				t.Fatalf("list len must match the level 0 nodes count. [cur:%v] [level0:%v]", cnt, list.Length)
			}
		}
	}
}

func TestBasicIntCRUD(t *testing.T) {
	var list *SkipList

	list = NewSkipList()

	list.Add(strconv.Itoa(10), IntToBytes(1))
	list.Add(strconv.Itoa(60), IntToBytes(2))
	list.Add(strconv.Itoa(30), IntToBytes(3))
	list.Add(strconv.Itoa(20), IntToBytes(4))
	list.Add(strconv.Itoa(90), IntToBytes(5))
	checkSanity(list, t)

	list.Add(strconv.Itoa(30), IntToBytes(9))
	checkSanity(list, t)

	list.Remove(strconv.Itoa(0))
	list.Remove(strconv.Itoa(20))
	checkSanity(list, t)

	v1, _ := list.Get(strconv.Itoa(10))
	v2, _ := list.Get(strconv.Itoa(60))
	v3, _ := list.Get(strconv.Itoa(30))
	v4, _ := list.Get(strconv.Itoa(20))
	v5, _ := list.Get(strconv.Itoa(90))
	v6, _ := list.Get(strconv.Itoa(0))

	if v1 == nil || !bytes.Equal(v1, IntToBytes(1)) {
		t.Fatal(`wrong "10" value (expected "1")`, v1)
	}

	if v2 == nil || !bytes.Equal(v2, IntToBytes(2)) {
		t.Fatal(`wrong "60" value (expected "2")`)
	}

	if v3 == nil || !bytes.Equal(v3, IntToBytes(9)) {
		t.Fatal(`wrong "30" value (expected "9")`)
	}

	if v4 != nil {
		t.Fatal(`found value for key "20", which should have been deleted`)
	}

	if v5 == nil || !bytes.Equal(v5, IntToBytes(5)) {
		t.Fatal(`wrong "90" value`)
	}

	if v6 != nil {
		t.Fatal(`found value for key "0", which should have been deleted`)
	}
}

func TestChangeLevel(t *testing.T) {
	var i int
	list := NewSkipList()

	if list.maxLevel != DefaultMaxLevel {
		t.Fatal("max level must equal default max value")
	}

	list = NewWithMaxLevel(4)
	if list.maxLevel != 4 {
		t.Fatal("wrong maxLevel (wanted 4)", list.maxLevel)
	}

	for i = 1; i <= 201; i++ {
		list.Add(strconv.Itoa(i), IntToBytes(i))
	}

	checkSanity(list, t)

	if list.Length != 201 {
		t.Fatal("wrong list length", list.Length)
	}

	for c := list.Front(); c != nil; c = c.Next() {
		if c.key != string(c.value) {
			t.Fatal("wrong list element value, key:", c.key, ", value:", string(c.value))
		}
	}
}

func TestMaxLevel(t *testing.T) {
	list := NewWithMaxLevel(DefaultMaxLevel + 1)
	list.Add(strconv.Itoa(0), []byte{})
}

func TestChangeProbability(t *testing.T) {
	list := NewSkipList()

	if list.probability != DefaultProbability {
		t.Fatal("new lists should have P value = DefaultProbability")
	}

	list.SetProbability(0.5)
	if list.probability != 0.5 {
		t.Fatal("failed to set new list probability value: expected 0.5, got", list.probability)
	}
}

func TestConcurrency(t *testing.T) {
	list := NewSkipList()

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		for i := 0; i < 100000; i++ {
			list.Add(strconv.Itoa(i), IntToBytes(i))
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < 100000; i++ {
			list.Get(strconv.Itoa(i))
		}
		wg.Done()
	}()

	wg.Wait()
	if list.Length != 100000 {
		t.Fail()
	}
}

func BenchmarkIncSet(b *testing.B) {
	b.ReportAllocs()
	list := NewSkipList()

	for i := 0; i < b.N; i++ {
		list.Add(strconv.Itoa(i), []byte{})
	}

	b.SetBytes(int64(b.N))
}

func BenchmarkIncGet(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		res, _ := benchList.Get(strconv.Itoa(i))
		if res == nil {
			b.Fatal("failed to Get an element that should exist")
		}
	}

	b.SetBytes(int64(b.N))
}

func BenchmarkDecSet(b *testing.B) {
	b.ReportAllocs()
	list := NewSkipList()

	for i := b.N; i > 0; i-- {
		list.Add(strconv.Itoa(i), []byte{})
	}

	b.SetBytes(int64(b.N))
}

func BenchmarkDecGet(b *testing.B) {
	b.ReportAllocs()
	for i := b.N; i > 0; i-- {
		res, _ := benchList.Get(strconv.Itoa(i))
		if res == nil {
			b.Fatal("failed to Get an element that should exist", i)
		}
	}

	b.SetBytes(int64(b.N))
}
