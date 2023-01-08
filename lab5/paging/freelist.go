package paging

import "fmt"

type freeList struct {
	freeList      []bool // tracks free physical frames
	numFreeFrames int    // number of free frames
}

// newFreeList creates a free list with space for numFrames frames.
func newFreeList(numFrames int) freeList {
	fl := freeList{numFreeFrames: numFrames}
	fl.freeList = make([]bool, numFrames)
	for i := range fl.freeList {
		fl.freeList[i] = true
	}
	return fl
}

// calculateNumFreeFrames returns the number of free frames in the free list.
func (fl freeList) calculateNumFreeFrames() int {
	n := 0
	for _, frame := range fl.freeList {
		if frame {
			n++
		}
	}
	return n
}

// removeFrames allocates memory by removing entries (indices) from the freeList.
// That is, sets them to false in the freeList and decrements numFreeFrames.
func (fl *freeList) removeFrames(entries []int) error {
	// will operate on a copy of the free list until it is determined no
	// errors occur, in order to "atomically" update it
	freeList := make([]bool, len(fl.freeList))
	copy(freeList, fl.freeList)

	// check the validity of each entry to be removed and modify the free list copy accordingly
	for _, entry := range entries {
		if entry >= len(freeList) || entry < 0 {
			return fmt.Errorf("failed to remove %d from free list: %w", entry, errIndexOutOfBounds)
		}
		if !fl.freeList[entry] {
			return errFreeListDuplicateOp
		}
		freeList[entry] = false
	}

	// finally update the MMU with the new state
	fl.freeList = freeList
	fl.numFreeFrames -= len(entries)
	return nil
}

// addFrames frees memory by adding entries (indices) to the freeList.
// That is, sets them to true in the freeList and increments numFreeFrames.
func (fl *freeList) addFrames(entries []int) error {
	// will operate on a copy of the free list until it is determined no
	// errors occur, in order to "atomically" update it
	freeList := make([]bool, len(fl.freeList))
	copy(freeList, fl.freeList)

	// check the validity of each entry to be added and modify the free list copy accordingly
	for _, entry := range entries {
		if entry >= len(fl.freeList) || entry < 0 {
			return fmt.Errorf("failed to add %d to free list: %w", entry, errIndexOutOfBounds)
		}
		if fl.freeList[entry] {
			return errFreeListDuplicateOp
		}
		freeList[entry] = true
	}

	// finally update the MMU with the new state
	fl.freeList = freeList
	fl.numFreeFrames += len(entries)
	return nil
}
