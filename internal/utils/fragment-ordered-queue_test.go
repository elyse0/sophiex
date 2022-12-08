package utils

import (
	"sophiex/internal/downloader/fragment"
	"testing"
)

func TestSimpleFragmentOrderedQueue(t *testing.T) {
	fragmentOrderedQueue := CreateFragmentOrderedQueue(3)

	fragmentOrderedQueue.Enqueue(fragment.FragmentResponse{
		Index:    0,
		Response: nil,
	})

	fragments, hasFinished := fragmentOrderedQueue.Dequeue()
	if len(fragments) != 1 {
		t.Errorf("Fragment dequeue should return one item")
	}
	if hasFinished {
		t.Errorf("Fragment dequeue should not have finished")
	}

	fragmentOrderedQueue.Enqueue(fragment.FragmentResponse{
		Index:    2,
		Response: nil,
	})

	fragments, hasFinished = fragmentOrderedQueue.Dequeue()
	if len(fragments) != 0 {
		t.Errorf("Fragment dequeue should not return items")
	}
	if hasFinished {
		t.Errorf("Fragment dequeue should not have finished")
	}

	fragmentOrderedQueue.Enqueue(fragment.FragmentResponse{
		Index:    1,
		Response: nil,
	})

	fragments, hasFinished = fragmentOrderedQueue.Dequeue()
	if len(fragments) != 2 {
		t.Errorf("Fragment dequeue should return two items")
	}

	if !hasFinished {
		t.Errorf("Fragment dequeue should have finished")
	}
}
