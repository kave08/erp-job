package fararavand

import "testing"

func TestFirstUnsyncedIndexReturnsFirstRecordAfterCheckpoint(t *testing.T) {
	t.Parallel()

	ids := []int{2, 4, 6, 9}
	index := firstUnsyncedIndex(len(ids), 4, func(i int) int {
		return ids[i]
	})

	if index != 2 {
		t.Fatalf("expected first unsynced index 2, got %d", index)
	}
}

func TestFirstUnsyncedIndexReturnsNegativeWhenBatchIsAlreadySynced(t *testing.T) {
	t.Parallel()

	ids := []int{2, 4, 6, 9}
	index := firstUnsyncedIndex(len(ids), 9, func(i int) int {
		return ids[i]
	})

	if index != -1 {
		t.Fatalf("expected fully synced batch to return -1, got %d", index)
	}
}
