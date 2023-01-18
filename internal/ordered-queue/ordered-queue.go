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

func (queue *OrderedQueue[T]) Enqueue(item OrderedItem[T]) {
	queue.items = append(queue.items, item)

	sort.Slice(queue.items, func(i, j int) bool {
		return queue.items[i].Index < queue.items[j].Index
	})
}

func (queue *OrderedQueue[T]) Dequeue() ([]OrderedItem[T], bool) {
	if len(queue.items) == 0 {
		return []OrderedItem[T]{}, false
	}

	if queue.items[0].Index != queue.current {
		return []OrderedItem[T]{}, false
	}

	cutIndex := 1
	for i := cutIndex; i < len(queue.items); i++ {
		if queue.items[i].Index != queue.items[i-1].Index+1 {
			break
		}

		cutIndex += 1
	}

	dequeueItems := queue.items[:cutIndex]
	lastItem := dequeueItems[len(dequeueItems)-1]
	queue.current = lastItem.Index + 1

	queue.items = queue.items[cutIndex:]

	return dequeueItems, lastItem.Index == (queue.itemsNumber - 1)
}
