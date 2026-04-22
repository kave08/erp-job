package transfer

import "testing"

func TestTrimAfterCheckpointReturnsFirstRecordAfterCheckpoint(t *testing.T) {
	t.Parallel()

	ids := []int{2, 4, 6, 9}
	trimmed := trimAfterCheckpoint(ids, 4, func(item int) int {
		return item
	})

	if len(trimmed) != 2 || trimmed[0] != 6 || trimmed[1] != 9 {
		t.Fatalf("unexpected trimmed result: %#v", trimmed)
	}
}

func TestTrimAfterCheckpointReturnsNilWhenBatchAlreadySynced(t *testing.T) {
	t.Parallel()

	ids := []int{2, 4, 6, 9}
	trimmed := trimAfterCheckpoint(ids, 9, func(item int) int {
		return item
	})

	if trimmed != nil {
		t.Fatalf("expected nil for fully synced batch, got %#v", trimmed)
	}
}
