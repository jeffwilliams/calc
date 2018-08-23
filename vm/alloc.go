package vm

type Alloced struct{}
type Free struct{}

// allocator is used for dynamic allocation on the data segment
type allocator struct {
	data []interface{}
	// base is the lowest index in data that can allocate; lower than that
	// is statically allocated before runtime
	//base int

	// free slots
	freeList []int
}

func (a *allocator) alloc() (slot int) {
	if len(a.freeList) > 0 {
		slot = a.freeList[len(a.freeList)-1]
		a.freeList = a.freeList[0 : len(a.freeList)-1]
		return
	}

	slot = len(a.freeList)
	a.data = append(a.data, Alloced{})
	return
}

func (a *allocator) free(slot int) {
	a.freeList = append(a.freeList, slot)
	a.data[slot] = Free{}
}
