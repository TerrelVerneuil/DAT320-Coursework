package paging

// MMU is the structure for the simulated memory management unit.
type MMU struct {
	frames    [][]byte           // contains memory content in form of frames[frameIndex][offset]
	freeList                     // tracks free physical frames
	processes map[int]*PageTable // contains page table for each process (key=pid)
	frameSize int
}

// OffsetLookupTable gives the bit mask corresponding to a virtual address's offset of length n,
// where n is the table index. This table can be used to find the offset mask needed to extract
// the offset from a virtual address. It supports up to 32-bit wide offset masks.
//
// OffsetLookupTable[0] --> 0000 ... 0000
// OffsetLookupTable[1] --> 0000 ... 0001
// OffsetLookupTable[2] --> 0000 ... 0011
// OffsetLookupTable[3] --> 0000 ... 0111
// OffsetLookupTable[8] --> 0000 ... 1111 1111
// etc.
var OffsetLookupTable = []int{
	// 0000, 0001, 0011, 0111, 1111, etc.
	0x0000000, 0x00000001, 0x00000003, 0x00000007,
	0x000000f, 0x0000001f, 0x0000003f, 0x0000007f,
	0x00000ff, 0x000001ff, 0x000003ff, 0x000007ff,
	0x0000fff, 0x00001fff, 0x00003fff, 0x00007fff,
	0x000ffff, 0x0001ffff, 0x0003ffff, 0x0007ffff,
	0x00fffff, 0x001fffff, 0x003fffff, 0x007fffff,
	0x0ffffff, 0x01ffffff, 0x03ffffff, 0x07ffffff,
	0xfffffff, 0x1fffffff, 0x3fffffff, 0x7fffffff, 0xffffffff,
}

// NewMMU creates a new MMU with a memory of memSize bytes.
// memSize should be >= 1 and a multiple of frameSize.
func NewMMU(memSize, frameSize int) *MMU {
	// TODO(student) Task 2: initialize the MMU object
	frame := memSize / frameSize
	frames := make([][]byte, frame) //capacity
	if memSize >= 1 && memSize%frameSize == 0 {
		for i := 0; i < len(frames); i++ {
			frames[i] = make([]byte, frameSize)
		}
	}
	return &MMU{
		frames:    frames,
		freeList:  newFreeList(frame),
		frameSize: frameSize,
		processes: make(map[int]*PageTable),
	}
}

// Alloc allocates n bytes of memory for process pid.
// The allocated memory is added to the process's page table.
// The process is given a page table if it doesn't already have one,
// unless an out of memory error occurred.
func (mmu *MMU) Alloc(pid, n int) error {
	// TODO(student) Task 2: implement memory allocation
	// Suggested approach:
	// - calculate #frames needed to allocate n bytes, error if not enough free frames
	// - if process pid has no page table, create one for it
	// - determine which frames to allocate to the process
	// - add the frames to the process's (identified by pid) page table and
	// - update the free list
	//fmt.Println(mmu.calculateNumFreeFrames())
	//var x int
	var needed int
	if n < 1 {
		return errNothingToAllocate
	}

	if n%mmu.frameSize == 0 {
		needed = n / mmu.frameSize
	} else {
		needed = (n / mmu.frameSize) + 1
	}
	if mmu.calculateNumFreeFrames() < needed {
		return errOutOfMemory
	}
	if mmu.processes[pid] == nil {
		mmu.processes[pid] = &PageTable{}
	}
	freeFrames, err := mmu.findFreeFrames(needed)
	if err != nil {
		return err
	}
	if err = mmu.removeFrames(freeFrames); err != nil {
		return err
	}
	mmu.processes[pid].Append(freeFrames)
	return nil
}

// Write writes content to the given process's address space starting at virtualAddress.
func (mmu *MMU) Write(pid, virtualAddress int, content []byte) error {
	// Suggested approach:
	// - check valid pid (must have a page table)
	// - translate the virtual address
	// - check if the memory must be extended in order to write the content
	//   from the given starting address
	// - attempt to allocate more memory if necessary to complete the write
	// - sequentially write content into the known-to-be-valid address space

	if mmu.processes[pid] == nil {
		return errInvalidProcess
	}
	vpn, offset, err := mmu.translateAndCheck(pid, virtualAddress) //valid virtual add
	if err != nil {
		return err
	}
	pt := mmu.processes[pid]
	pleft := pt.Len() - vpn //
	bytesl := (pleft * mmu.frameSize) - offset

	if len(content) > bytesl {
		x := len(content) - bytesl
		err := mmu.Alloc(pid, x)
		if err != nil {
			return err
		}
	}
	for i := range content {
		pframe, _ := pt.Lookup(vpn)
		mmu.frames[pframe][offset] = content[i]
		offset++
		if offset >= mmu.frameSize {
			vpn++
			offset = 0
		}
	}
	return nil
}

// Read returns content of size n bytes from the given process's address space starting at virtualAddress.
func (mmu *MMU) Read(pid, virtualAddress, n int) (content []byte, err error) {
	// TODO(student) Task 3: implement reading
	// Suggested approach:
	// - check valid pid (must have a page table)
	// - translate the virtual address
	// - read and return the requested memory content
	if n < 1 { //there is nothing to read if the read bytes are less than 1
		return content, errNothingToRead
	}
	if mmu.processes[pid] == nil { //not a valid pid
		return content, errInvalidProcess
	}
	vpn, offset, err := mmu.translateAndCheck(pid, virtualAddress) //valid virtual add
	if err != nil {
		return content, err
	}
	pt := mmu.processes[pid]
	pleft := pt.Len() - vpn
	bytesl := (pleft * mmu.frameSize) - offset

	if n > bytesl { //if the bytes left is less than n than there is no more
		//free memory so we return an error
		return content, errOutOfMemory
	}
	//p
	for i := 0; i < n; i++ {
		pframe, _ := pt.Lookup(vpn) //we get the mapping of the virtual page
		//number from the vpn we got from translate and check
		content = append(content, mmu.frames[pframe][offset])
		offset++ //increase the offset
		if offset >= mmu.frameSize {
			vpn++      //increase the vpn
			offset = 0 //we reset to 0 because we reached the end of the frame
			//so we move on to the next adding content to the frame
		}
	}

	//mmu.frames is the memory content of the frame
	//return content of size n bytes
	return content, nil
}

// Free is called by a process's Free() function to free some of its allocated memory.
func (mmu *MMU) Free(pid, n int) error {
	// TODO(student) Task 4: implement freeing of memory
	// Suggested approach:
	// - check valid pid (must have a page table)
	// - check if there are at least n entries in the page table of pid
	// - free n pages
	// - set all the bytes in the freed memory to the value 0
	// - re-add the freed frames to the free list

	// - check valid pid (must have a page table) *  //if process at pagetable exists
	if mmu.processes[pid] == nil {
		return errInvalidProcess
	}
	// - check if there are at least n entries in the page table of pid

	// - free n pages
	freed, err := mmu.processes[pid].Free(n)
	if err != nil {
		return err
	}
	for _, phyin := range freed {
		for offset := range mmu.frames[phyin] {
			mmu.frames[phyin][offset] = 0
		}
	}
	//set all bytes in the freed memory to 0
	err = mmu.freeList.addFrames(freed) //add back the freed frames
	if err != nil {
		return err
	}
	return nil
}

// extract returns the virtual page number and offset for the given virtual address,
// and the number of bits in the offset n.
func extract(virtualAddress, n int) (vpn, offset int) {
	// TODO(student) Implement virtual address translation as described in
	// the Virtual Addresses section of the README.
	offset = OffsetLookupTable[n] & virtualAddress
	vpn = virtualAddress >> n
	return vpn, offset
}

// translateAndCheck returns the virtual page number and offset for the given virtual address.
// If the virtual address is invalid for process pid, an error is returned.
func (mmu *MMU) translateAndCheck(pid, virtualAddress int) (vpn, offset int, err error) {
	// TODO(student) Implement virtual address translation as described in
	// the Virtual Addresses section of the README.
	// The procedure is described in detail in Chapter 18.1 of the textbook. *
	// It is expected that this method calls the extract function above *
	// to compute the VPN and offset to be returned from this function after
	// checking that the process has access to the returned VPN.
	// You might also find the provided log2 function useful to calculate one
	// of the inputs to the extract function.
	//_, err = mmu.processes[pid].Lookup(vpn) //has access
	//if err != nil {
	//	return 0, 0, err
	//}
	n := log2(mmu.frameSize) // n is is given by log framesize
	vpn, offset = extract(virtualAddress, n)
	_, err = mmu.processes[pid].Lookup(vpn)
	if err != nil {
		return 0, 0, err
	}
	return vpn, offset, nil
}

// log2 calculates m given n = 2^m.
func log2(n int) int {
	exp := 0
	for {
		if n%2 == 0 && n > 0 {
			exp++
		} else {
			return exp
		}
		n /= 2
	}
}
