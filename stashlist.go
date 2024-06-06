package stashlist

import (
	"math"
	"math/rand"
	"time"
)

const (
	DefaultMaxLevel    int     = 18
	DefaultProbability float64 = 1 / math.E
)

type elementNode struct {
	next    []*Element
	level   int
	visited bool
}

type Element struct {
	elementNode
	key   string
	value []byte
}

// Next returns the following Element or nil if we're at the end of the list.
// Only operates on the bottom level of the skip list (a fully linked list).
func (element *Element) Next() *Element {
	return element.next[0]
}

type StashList struct {
	elementNode
	maxLevel       int
	Length         int
	randSource     rand.Source
	probability    float64
	probTable      []float64
	prevNodesCache []*elementNode
}

// Front returns the head node of the list.
func (list *StashList) Front() *Element {
	return list.next[0]
}

// Add inserts a value in the list with the specified key, ordered by the key.
// If the key exists, it updates the value in the existing node.
// Returns a pointer to the new element.
func (list *StashList) Add(key string, value []byte) {
	var element *Element
	prevs := list.getPrevElementNodes(key)

	if element = prevs[0].next[0]; element != nil && element.key <= key {
		if element.visited == false {
			element.visited = true
		} else {
			// Promote
			level := element.level
			if level < list.maxLevel && prevs[level] != &list.elementNode {
				element.next[level] = prevs[level].next[level]
				prevs[level].next[level] = element
				if prevs[level].visited == true {
					prevs[level].visited = false
				}

				level = level + 1
				element.level = level
			}
		}
		element.value = value
		return
	}

	level := list.randLevel()
	element = &Element{
		elementNode: elementNode{
			next:    make([]*Element, list.maxLevel),
			level:   level,
			visited: false,
		},
		key:   key,
		value: value,
	}
	if level == 1 {
		element.visited = true
	}

	for i := 0; i < level; i++ {
		element.next[i] = prevs[i].next[i]
		prevs[i].next[i] = element
	}

	list.Length++
}

// Get finds an element by key. It returns element pointer if found, nil if not found.
func (list *StashList) Get(key string) ([]byte, bool) {
	var prev = &list.elementNode
	var next *Element

	for i := list.maxLevel - 1; i >= 0; i-- {
		next = prev.next[i]

		for next != nil && key > next.key {
			prev = &next.elementNode
			next = next.next[i]
		}
	}

	if next != nil && next.key <= key {
		if next.visited == false {
			next.visited = true
		}
		return next.value, true
	}

	return nil, false
}

// Remove deletes an element from the list.
// Returns removed element pointer if found, nil if not found.
func (list *StashList) Remove(key string) *Element {
	prevs := list.getPrevElementNodes(key)

	// found the element, remove it
	if element := prevs[0].next[0]; element != nil && element.key <= key {
		for k, v := range element.next {
			prevs[k].next[k] = v
		}

		list.Length--
		return element
	}

	return nil
}

// getPrevElementNodes is the private search mechanism that other functions use.
// Finds the previous nodes on each level relative to the current Element and
// caches them. This approach is similar to a "search finger" as described by Pugh:
// http://citeseerx.ist.psu.edu/viewdoc/summary?doi=10.1.1.17.524
func (list *StashList) getPrevElementNodes(key string) []*elementNode {
	var prev = &list.elementNode
	var next *Element
	var before *elementNode

	prevs := list.prevNodesCache

	for i := list.maxLevel - 1; i >= 0; i-- {
		next = prev.next[i]

		for next != nil && key > next.key {
			before = prev
			prev = &next.elementNode
			next = next.next[i]
			if i > 0 && next != nil && key == next.key && prev.visited == false {
				// Demote
				before.next[i] = next
				prev.next[i] = nil
				prev.level = prev.level - 1
				break
			}
			// TODO: flush unvisited items
			// if i == 0 {}
		}

		prevs[i] = prev
	}

	return prevs
}

// SetProbability changes the current P value of the list.
// It doesn't alter any existing data, only changes how future insert heights are calculated.
func (list *StashList) SetProbability(newProbability float64) {
	list.probability = newProbability
	list.probTable = probabilityTable(list.probability, list.maxLevel)
}

func (list *StashList) randLevel() (level int) {
	// Our random number source only has Int63(), so we have to produce a float64 from it
	// Reference: https://golang.org/src/math/rand/rand.go#L150
	r := float64(list.randSource.Int63()) / (1 << 63)

	level = 1
	for level < list.maxLevel && r < list.probTable[level] {
		level++
	}
	return
}

// probabilityTable calculates in advance the probability of a new node having a given level.
// probability is in [0, 1], MaxLevel is (0, 64]
// Returns a table of floating point probabilities that each level should be included during an insert.
func probabilityTable(probability float64, MaxLevel int) (table []float64) {
	for i := 1; i <= MaxLevel; i++ {
		prob := math.Pow(probability, float64(i-1))
		table = append(table, prob)
	}
	return table
}

// NewWithMaxLevel creates a new skip list with MaxLevel set to the provided number.
// maxLevel has to be int(math.Ceil(math.Log(N))) for DefaultProbability (where N is an upper bound on the
// number of elements in a skip list). See http://citeseerx.ist.psu.edu/viewdoc/summary?doi=10.1.1.17.524
// Returns a pointer to the new list.
func NewWithMaxLevel(maxLevel int) *StashList {
	if maxLevel < 1 || maxLevel > 64 {
		panic("maxLevel for a StashList must be a positive integer <= 64")
	}

	return &StashList{
		elementNode:    elementNode{next: make([]*Element, maxLevel), level: maxLevel},
		prevNodesCache: make([]*elementNode, maxLevel),
		maxLevel:       maxLevel,
		randSource:     rand.New(rand.NewSource(time.Now().UnixNano())),
		probability:    DefaultProbability,
		probTable:      probabilityTable(DefaultProbability, maxLevel),
	}
}

// NewStashList creates a new skip list with default parameters. Returns a pointer to the new list.
func NewStashList() *StashList {
	return NewWithMaxLevel(DefaultMaxLevel)
}
