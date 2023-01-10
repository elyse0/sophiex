package ordered_queue

import (
	"sort"
)

type OrderedItem[T any] struct {
	Index   int
	Payload T
}

type OrderedQueue[T any] struct {
	itemsNumber int
	current     int
	items       []OrderedItem[T]
}

func CreateOrderedQueue[T any](itemsNumber int) *OrderedQueue[T] {
	return &OrderedQueue[T]{
		itemsNumber: itemsNumber,
		current:     0,
		items:       []OrderedItem[T]{},
	}
}

func (slice *OrderedQueue[T]) Enqueue(item OrderedItem[T]) {
	slice.items = append(slice.items, item)

	sort.Slice(slice.items, func(i, j int) bool {
		return slice.items[i].Index < slice.items[j].Index
	})
}

func (slice *OrderedQueue[T]) Dequeue() ([]OrderedItem[T], bool) {
	if len(slice.items) == 0 {
		return []OrderedItem[T]{}, false
	}

	if slice.items[0].Index != slice.current {
		return []OrderedItem[T]{}, false
	}

	cutIndex := 1
	for i := cutIndex; i < len(slice.items); i++ {
		if slice.items[i].Index != slice.items[i-1].Index+1 {
			break
		}

		cutIndex += 1
	}

	dequeueItems := slice.items[:cutIndex]
	lastItem := dequeueItems[len(dequeueItems)-1]
	slice.current = lastItem.Index + 1

	slice.items = slice.items[cutIndex:]

	return dequeueItems, lastItem.Index == (slice.itemsNumber - 1)
}
