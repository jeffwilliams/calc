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

func newAllocator(data []interface{}) allocator {
	return allocator{data: data, freeList: make([]int, 0, 100)}
}

func (a *allocator) alloc() (slot int) {
	if len(a.freeList) > 0 {
		slot = a.freeList[len(a.freeList)-1]
		a.data[slot] = Alloced{}
		a.freeList = a.freeList[0 : len(a.freeList)-1]
		return
	}

	slot = len(a.data)
	a.data = append(a.data, Alloced{})
	return
}

func (a *allocator) free(slot int) {
	if _, ok := a.data[slot].(Free); ok {
		// already free
		return
	}

	a.freeList = append(a.freeList, slot)
	a.data[slot] = Free{}
}
