package paging

// NoEntry is produced when no entry matching a request exists
const NoEntry = -1

// PageTable is a per-process data structure which holds translations from virtual page numbers to physical frame numbers
type PageTable struct {
	frameIndices []int // maps virtual page number (index) to physical frame number (content)
}

// Append adds pages to a page table
func (pt *PageTable) Append(pages []int) {
	pt.frameIndices = append(pt.frameIndices, pages...)
}

// Free removes the n last pages from the page table and returns the removed entries
func (pt *PageTable) Free(n int) ([]int, error) {
	// TODO(student) Implement free functionality for the pagetable
	var removed []int
	if n < 1 || n > pt.Len() {
		return []int{}, errFreeOutOfBounds
	}
	if len(pt.frameIndices) > 0 && n <= pt.Len() { //slice is not empty
		removed = pt.frameIndices[len(pt.frameIndices)-n:]
		remained := pt.frameIndices[:len(pt.frameIndices)-n] //re
		pt.frameIndices = remained
	} else {
		return []int{}, errInvalidProcess
	}

	return removed, nil
}

// Lookup returns the mapping of a virtual page number to a physical frame number, or an error if it does not exist.
func (pt *PageTable) Lookup(virtualPageNum int) (frameIndex int, err error) {
	// TODO(student) Implement lookup functionality for the pagetable
	if virtualPageNum >= pt.Len() || virtualPageNum < 0 {
		return NoEntry, errIndexOutOfBounds
	}
	frameIndex = pt.frameIndices[virtualPageNum]
	return frameIndex, nil
}

// Len returns the length of the page table
func (pt *PageTable) Len() int {
	return len(pt.frameIndices)
}
