package paging

import "fmt"

// PrintMemory will print the contents of the memory for debugging purposes.
func (mmu *MMU) PrintMemory() {
	frameSize := len(mmu.frames[0])
	for i, frame := range mmu.frames {
		fmt.Printf("[%s: ", fmt.Sprintf("0x%x", i*frameSize))
		if mmu.freeList.freeList[i] {
			fmt.Print("FREE]\n")
		} else {
			fmt.Print("BUSY]\n")
		}
		for _, cell := range frame {
			fmt.Printf("> %08b\n", cell)
		}
	}
	fmt.Println("------------")
}

// setMemoryContent sets the memory content (mmu.frames) to a certain state.
// It is used in testing. If your implementation requires additional actions
// to be done, you can define them here.
func (mmu *MMU) setMemoryContent(frames [][]byte) {
	mmu.frames = frames
}

// setFreeList sets the free list to a certain state.
// It is used in testing. If your implementation requires additional actions
// to be done, you can define them here.
func (mmu *MMU) setFreeList(freeList []bool) {
	mmu.freeList.freeList = freeList
	mmu.freeList.numFreeFrames = mmu.calculateNumFreeFrames()
}

// setProcesses sets the state of multiple processes.
// It is used in testing. If your implementation requires additional actions
// to be done, you can define them here.
func (mmu *MMU) setProcesses(processes map[int]*PageTable) {
	for i, p := range processes {
		mmu.setProcess(i, p)
	}
}

// setProcesses sets the state of a single process.
// It is used in testing. If your implementation requires additional actions
// to be done, you can define them here.
func (mmu *MMU) setProcess(pid int, process *PageTable) {
	copyProcess := *process
	mmu.processes[pid] = &copyProcess
}
