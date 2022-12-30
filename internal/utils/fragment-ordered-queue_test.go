package utils

import (
	"net/http"
	"testing"
)

func TestSimpleFragmentOrderedQueue(t *testing.T) {
	fragmentOrderedQueue := CreateFragmentOrderedQueue[*http.Response](3)

	fragmentOrderedQueue.Enqueue(OrderedFragment[*http.Response]{
		Index:   0,
		Payload: nil,
	})

	fragments, hasFinished := fragmentOrderedQueue.Dequeue()
	if len(fragments) != 1 {
		t.Errorf("Fragment dequeue should return one item")
	}
	if hasFinished {
		t.Errorf("Fragment dequeue should not have finished")
	}

	fragmentOrderedQueue.Enqueue(OrderedFragment[*http.Response]{
		Index:   2,
		Payload: nil,
	})

	fragments, hasFinished = fragmentOrderedQueue.Dequeue()
	if len(fragments) != 0 {
		t.Errorf("Fragment dequeue should not return items")
	}
	if hasFinished {
		t.Errorf("Fragment dequeue should not have finished")
	}

	fragmentOrderedQueue.Enqueue(OrderedFragment[*http.Response]{
		Index:   1,
		Payload: nil,
	})

	fragments, hasFinished = fragmentOrderedQueue.Dequeue()
	if len(fragments) != 2 {
		t.Errorf("Fragment dequeue should return two items")
	}

	if !hasFinished {
		t.Errorf("Fragment dequeue should have finished")
	}
}
