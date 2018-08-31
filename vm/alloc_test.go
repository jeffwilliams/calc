package vm

import "testing"

func TestEmpty(t *testing.T) {
	data := []interface{}{}
	alloc := newAllocator(data)

	slt := alloc.alloc()

	if slt != 0 {
		t.Fatalf("wrong slot")
	}

	slt2 := alloc.alloc()
	if slt2 != 1 {
		t.Fatalf("wrong slot")
	}

	alloc.free(slt2)
	if len(alloc.freeList) == 0 || alloc.freeList[0] != 1 {
		t.Fatalf("wrong free list")
	}

	slt2 = alloc.alloc()
	if slt2 != 1 {
		t.Fatalf("wrong slot")
	}

	if len(alloc.freeList) != 0 {
		t.Fatalf("wrong free list")
	}

	alloc.free(slt2)
	t.Logf("freeing %d", slt2)
	alloc.free(slt)
	t.Logf("freeing %d", slt)

	if len(alloc.freeList) != 2 {
		t.Fatalf("wrong free list. size = %d %v", len(alloc.freeList), alloc.freeList)
	}

	// dbl free
	alloc.free(slt)
	if len(alloc.freeList) != 2 {
		t.Fatalf("wrong free list")
	}
}

func TestNonEmpty(t *testing.T) {
	data := []interface{}{4, 5}
	alloc := newAllocator(data)

	slt := alloc.alloc()

	if slt != 2 {
		t.Fatalf("wrong slot")
	}
}
