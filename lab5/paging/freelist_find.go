package paging

// findFreeFrames returns indices for n free frames.
// If there are not enough free frames available, an error is returned.
func (fl *freeList) findFreeFrames(n int) ([]int, error) {
	var x []int
	for i := 0; i < len(fl.freeList); i++ {
		if fl.freeList[i] == true {
			x = append(x, i) //only 2
		}
		if len(x) == n {
			return x, nil
		}
	}
	//there isn't enough frames
	return x, errNothingToAllocate
}
