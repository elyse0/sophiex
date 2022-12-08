package utils

import (
	"sophiex/internal/downloader/fragment"
	"sort"
)

type FragmentOrderedQueue struct {
	itemsNumber int
	current     int
	items       []fragment.FragmentResponse
}

func CreateFragmentOrderedQueue(itemsNumber int) *FragmentOrderedQueue {
	return &FragmentOrderedQueue{
		itemsNumber: itemsNumber,
		current:     0,
		items:       []fragment.FragmentResponse{},
	}
}

func (slice *FragmentOrderedQueue) Enqueue(item fragment.FragmentResponse) {
	slice.items = append(slice.items, item)

	sort.Slice(slice.items, func(i, j int) bool {
		return slice.items[i].Index < slice.items[j].Index
	})
}

func (slice *FragmentOrderedQueue) Dequeue() ([]fragment.FragmentResponse, bool) {
	if len(slice.items) == 0 {
		return []fragment.FragmentResponse{}, false
	}

	if slice.items[0].Index != slice.current {
		return []fragment.FragmentResponse{}, false
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
