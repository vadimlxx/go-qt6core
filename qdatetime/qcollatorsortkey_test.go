package qdatetime

import (
	"runtime"
	"testing"
)

// newTestSortKeys returns default ("C"-like) sort keys for two ASCII strings.
func newTestSortKeys(t *testing.T) (apple, banana *QCollatorSortKey) {
	t.Helper()
	apple = DefaultSortKeyString("apple")
	banana = DefaultSortKeyString("banana")
	if apple == nil || banana == nil {
		t.Fatal("DefaultSortKeyString returned nil")
	}
	t.Cleanup(apple.Close)
	t.Cleanup(banana.Close)
	return apple, banana
}

func TestQCollatorSortKeyCompare(t *testing.T) {
	apple, banana := newTestSortKeys(t)

	if got := apple.Compare(banana); got >= 0 {
		t.Fatalf("apple.Compare(banana) = %d, want < 0", got)
	}
	if got := banana.Compare(apple); got <= 0 {
		t.Fatalf("banana.Compare(apple) = %d, want > 0", got)
	}
	if got := apple.Compare(apple); got != 0 {
		t.Fatalf("apple.Compare(apple) = %d, want 0", got)
	}
}

func TestQCollatorSortKeyLess(t *testing.T) {
	apple, banana := newTestSortKeys(t)

	if !apple.Less(banana) {
		t.Fatal("apple.Less(banana) = false, want true")
	}
	if banana.Less(apple) {
		t.Fatal("banana.Less(apple) = true, want false")
	}
	if apple.Less(apple) {
		t.Fatal("apple.Less(apple) = true, want false")
	}
}

func TestQCollatorSortKeyCopy(t *testing.T) {
	apple, banana := newTestSortKeys(t)

	cp := NewQCollatorSortKeyCopy(apple)
	if cp == nil {
		t.Fatal("NewQCollatorSortKeyCopy returned nil")
	}
	defer cp.Close()
	if cp.Compare(apple) != 0 {
		t.Fatal("copy must compare equal to the original")
	}
	if cp.Compare(banana) >= 0 {
		t.Fatal("copy must still sort before banana")
	}

	generic := NewQCollatorSortKey(apple)
	if generic == nil {
		t.Fatal("NewQCollatorSortKey(key) returned nil")
	}
	defer generic.Close()
	if generic.Compare(apple) != 0 {
		t.Fatal("generic copy must compare equal to the original")
	}

	clone := apple.Clone()
	if clone == nil {
		t.Fatal("Clone returned nil")
	}
	defer clone.Close()
	if clone.Compare(apple) != 0 {
		t.Fatal("Clone must compare equal to the original")
	}
}

func TestQCollatorSortKeySwap(t *testing.T) {
	apple, banana := newTestSortKeys(t)

	apple.Swap(banana)
	// After the swap, apple holds banana's key and banana holds apple's key.
	if apple.Less(banana) {
		t.Fatal("after Swap, apple key must not sort before banana key")
	}
	if !banana.Less(apple) {
		t.Fatal("after Swap, banana key must sort before apple key")
	}
}

func TestQCollatorSortKeyGenericConstructorPanics(t *testing.T) {
	for _, arg := range [][]any{nil, []any{"not a key"}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Fatalf("NewQCollatorSortKey(%v) must panic", arg)
				}
			}()
			NewQCollatorSortKey(arg...)
		}()
	}
}

func TestQCollatorSortKeyClose(t *testing.T) {
	k := DefaultSortKeyString("x")
	if k == nil {
		t.Fatal("DefaultSortKeyString returned nil")
	}
	k.Close()
	k.Close() // idempotent

	y := DefaultSortKeyString("y")
	if y == nil {
		t.Fatal("DefaultSortKeyString returned nil")
	}
	defer y.Close()

	defer func() {
		if recover() == nil {
			t.Fatal("Compare on a closed sort key must panic")
		}
	}()
	k.Compare(y)
}

func TestQCollatorSortKeyCleanup(t *testing.T) {
	for i := 0; i < 32; i++ {
		k := DefaultSortKeyString("key")
		if k == nil {
			t.Fatalf("iteration %d: nil sort key", i)
		}
	}
	runtime.GC()
	runtime.GC()
}
