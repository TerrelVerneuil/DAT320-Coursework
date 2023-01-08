package paging

import (
	"errors"
	"fmt"

	"github.com/google/go-cmp/cmp"
)

// -------------
// | MMU tests |
// -------------

var (
	cmpOptPageTable         = cmp.AllowUnexported(PageTable{})
	errAllocInvalidInput    = errors.New("cannot allocate less than 1 byte")
	errAllocNotEnoughFrames = errors.New("not enough free frames to allocate what the process requested")
	errWriteOutOfBounds     = errors.New("tried to write to unallocated or non-existant address")
	errReadInvalidAddr      = errors.New("tried to read from an unallocated or non-existant virtual address")
	errReadOutOfBounds      = errors.New("tried to read outside of the bounds of the process' virtual address space")
	errFreeTooManyPages     = errors.New("tried to free more pages than were allocated to the process")
)

var newMMUTests = []struct {
	memSize, frameSize int // parameters to NewMMU
	wantFrames         [][]byte
	wantFreeList       []bool
	desc               string
}{
	{
		memSize: 1, frameSize: 1,
		wantFrames:   [][]byte{{0}},
		wantFreeList: []bool{true},
		desc:         "create MMU with 1 frame of size 1",
	},
	{
		memSize: 2, frameSize: 1,
		wantFrames:   [][]byte{{0}, {0}},
		wantFreeList: []bool{true, true},
		desc:         "create MMU with 2 frames of size 1",
	},
	{
		memSize: 2, frameSize: 2,
		wantFrames:   [][]byte{{0, 0}},
		wantFreeList: []bool{true},
		desc:         "create MMU with 1 frame of size 2",
	},
	{
		memSize: 4, frameSize: 1,
		wantFrames:   [][]byte{{0}, {0}, {0}, {0}},
		wantFreeList: []bool{true, true, true, true},
		desc:         "create MMU with 4 frames of size 1",
	},
	{
		memSize: 16, frameSize: 4,
		wantFrames: [][]byte{
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		wantFreeList: []bool{true, true, true, true},
		desc:         "create MMU with 4 frames of size 4",
	},
	{
		memSize: 64, frameSize: 8,
		wantFrames: [][]byte{
			{0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0},
		},
		wantFreeList: []bool{true, true, true, true, true, true, true, true},
		desc:         "create MMU with 8 frames of size 8",
	},
}

var AllocTests = []struct {
	in                 int        // number of bytes to Alloc for pid=0
	pageTable          *PageTable // expected page table for pid=0 after the operation
	err                error      // expected error
	freeList           []bool     // free list before the operation
	memSize, frameSize int        // parameters to NewMMU
	desc               string
}{
	{
		in:        0,
		pageTable: nil,
		err:       errAllocInvalidInput,
		freeList:  []bool{true},
		memSize:   1, frameSize: 1,
		desc: "Tried to request '0' bytes - must request at least 1 byte",
	},
	{
		in:        -1,
		pageTable: nil,
		err:       errAllocInvalidInput,
		freeList:  []bool{true},
		memSize:   1, frameSize: 1,
		desc: "Tried to request '-1' bytes - must request at least 1 byte",
	},
	{
		in:        1,
		pageTable: nil,
		err:       errAllocNotEnoughFrames,
		freeList:  []bool{false},
		memSize:   1, frameSize: 1,
		desc: "Valid request, but system is out of memory",
	},
	{
		in:        1,
		pageTable: &PageTable{[]int{0}},
		err:       nil,
		freeList:  []bool{true},
		memSize:   1, frameSize: 1,
		desc: "Process requests and is allocated 1 frame",
	},
	{
		in:        1,
		pageTable: &PageTable{[]int{1}},
		err:       nil,
		freeList:  []bool{false, true},
		memSize:   2, frameSize: 1,
		desc: "Process requests 1 frame and is allocated the 2nd frame since the 1st frame is occupied",
	},
	{
		in:        1,
		pageTable: &PageTable{[]int{0}},
		err:       nil,
		freeList:  []bool{true, true},
		memSize:   2, frameSize: 1,
		desc: "Process requests 1 frame and is allocated the 1st out of 2 consecutive free frames",
	},
	{
		in:        2,
		pageTable: &PageTable{[]int{0, 1}},
		err:       nil,
		freeList:  []bool{true, true},
		memSize:   2, frameSize: 1,
		desc: "Process requests 2 frames and is allocated both free frames",
	},
	{
		in:        3,
		pageTable: nil,
		err:       errAllocNotEnoughFrames,
		freeList:  []bool{true, false},
		memSize:   4, frameSize: 2,
		desc: "Process requests 1 frame and is allocated the 1st out of 2 consecutive free frames",
	},
	{
		in:        4,
		pageTable: &PageTable{[]int{0, 1}},
		err:       nil,
		freeList:  []bool{true, true},
		memSize:   4, frameSize: 2,
		desc: "Process requests 4 bytes and is allocated 2 frames since frame size is 2 bytes",
	},
	{
		in:        3,
		pageTable: &PageTable{[]int{1, 2}},
		err:       nil,
		freeList:  []bool{false, true, true},
		memSize:   6, frameSize: 2,
		desc: "Process requests 3 bytes and is allocated 2 frames (3 bytes rounded up to 4 bytes to fit frame sizes) since frame size is 2 bytes.\nThe process is allocated frames 1 and 2 since frame 0 is occupied.",
	},
	{
		in:        5,
		pageTable: &PageTable{[]int{0, 1, 2}},
		err:       nil,
		freeList:  []bool{true, true, true},
		memSize:   6, frameSize: 2,
		desc: "Process requests 5 bytes and is allocated 3 frames (5 bytes rounded up to 6 bytes rounded to fit frame sizes) since frame size is 2 bytes",
	},
	{
		in:        6,
		pageTable: &PageTable{[]int{0, 1, 2}},
		err:       nil,
		freeList:  []bool{true, true, true},
		memSize:   6, frameSize: 2,
		desc: "Process requests 6 bytes and is allocated 3 frames since frame size is 2 bytes",
	},
	{
		in:        16,
		pageTable: &PageTable{[]int{1, 2, 4, 5}},
		err:       nil,
		freeList:  []bool{false, true, true, false, true, true, false, true},
		memSize:   32, frameSize: 4,
		desc: "Process requests 16 bytes and is allocated 4 frames (spread in memory) since frame size is 4 bytes",
	},
	{
		in:        17,
		pageTable: &PageTable{[]int{1, 2, 4, 5, 7}},
		err:       nil,
		freeList:  []bool{false, true, true, false, true, true, false, true},
		memSize:   32, frameSize: 4,
		desc: "Process requests 17 bytes and is allocated 5 frames (17 bytes rounded up to 20 bytes to fit frame size) (spread in memory) since frame size is 4 bytes",
	},
}

type TAllocMultipleOperation struct {
	in, pid       int        // in - bytes to request allocated by pid
	wantPageTable *PageTable // state of page table for process pid, after the operation
	wantFreeList  []bool     // state of the free list after the operation
	wantError     error      // error expected from the operation
	desc          string     // what this operation is testing
}

var AllocMultipleTests = []struct {
	operations         []TAllocMultipleOperation // several Alloc operations to be performed in sequence
	memSize, frameSize int                       // parameters to NewMMU
	freeListState      []bool                    // state of MMU.freeList before the test
}{
	{
		operations: []TAllocMultipleOperation{
			{
				in: 1, pid: 0,
				wantPageTable: nil,
				wantFreeList:  []bool{false},
				wantError:     errAllocNotEnoughFrames,
				desc:          "no memory available -> error",
			},
		},
		memSize:       1,
		frameSize:     1,
		freeListState: []bool{false},
	},
	{
		operations: []TAllocMultipleOperation{
			{
				in: 1, pid: 0,
				wantPageTable: &PageTable{[]int{0}},
				wantFreeList:  []bool{false},
				wantError:     nil,
				desc:          "memory available -> allocate",
			},
		},
		memSize:       1,
		frameSize:     1,
		freeListState: []bool{true},
	},
	{
		operations: []TAllocMultipleOperation{
			{
				in: 1, pid: 0,
				wantPageTable: &PageTable{[]int{0}},
				wantFreeList:  []bool{false},
				wantError:     nil,
				desc:          "valid allocation",
			},
			{
				in: 1, pid: 0,
				wantPageTable: &PageTable{[]int{0}},
				wantFreeList:  []bool{false},
				wantError:     errAllocNotEnoughFrames,
				desc:          "out of memory -> error and no allocation",
			},
		},
		memSize:       1,
		frameSize:     1,
		freeListState: []bool{true},
	},
	{
		operations: []TAllocMultipleOperation{
			{
				in: 1, pid: 0,
				wantPageTable: &PageTable{[]int{0}},
				wantFreeList:  []bool{false, true},
				wantError:     nil,
				desc:          "allocate several times -> add to page table",
			},
			{
				in: 1, pid: 0,
				wantPageTable: &PageTable{[]int{0, 1}},
				wantFreeList:  []bool{false, false},
				wantError:     nil,
				desc:          "allocate several times -> add to page table",
			},
		},
		memSize:       2,
		frameSize:     1,
		freeListState: []bool{true, true},
	},
	{
		operations: []TAllocMultipleOperation{
			{
				in: 1, pid: 0,
				wantPageTable: &PageTable{[]int{1}},
				wantFreeList:  []bool{false, false},
				wantError:     nil,
				desc:          "allocate 2nd frame (1st frame is not free) -> page table points to 2nd frame correctly",
			},
		},
		memSize:       2,
		frameSize:     1,
		freeListState: []bool{false, true},
	},
	{
		operations: []TAllocMultipleOperation{
			{
				in: 2, pid: 0,
				wantPageTable: &PageTable{[]int{0, 1}},
				wantFreeList:  []bool{false, false},
				wantError:     nil,
				desc:          "allocate more than 1 frame",
			},
		},
		memSize:       2,
		frameSize:     1,
		freeListState: []bool{true, true},
	},
	{
		operations: []TAllocMultipleOperation{
			{
				in: 2, pid: 0,
				wantPageTable: &PageTable{[]int{1, 3}},
				wantFreeList:  []bool{false, false, false, false},
				wantError:     nil,
				desc:          "allocate 2 frames that are not contiguous in memory layout",
			},
		},
		memSize:       4,
		frameSize:     1,
		freeListState: []bool{false, true, false, true},
	},
	{
		operations: []TAllocMultipleOperation{
			{
				in: 5, pid: 0,
				wantPageTable: &PageTable{[]int{0}},
				wantFreeList:  []bool{false},
				wantError:     nil,
				desc:          "round up requested bytes to nearest multiple of frame size",
			},
		},
		memSize:       8,
		frameSize:     8,
		freeListState: []bool{true},
	},
	{
		operations: []TAllocMultipleOperation{
			{
				in: 9, pid: 0,
				wantPageTable: &PageTable{[]int{0, 1}},
				wantFreeList:  []bool{false, false},
				wantError:     nil,
				desc:          "round up requested bytes to nearest multiple of frame size",
			},
		},
		memSize:       16,
		frameSize:     8,
		freeListState: []bool{true, true},
	},
	{
		operations: []TAllocMultipleOperation{
			{
				in: 9, pid: 0,
				wantPageTable: nil,
				wantFreeList:  []bool{true, false},
				wantError:     errAllocNotEnoughFrames,
				desc:          "only 8 bytes free, while 9 were requested -> error and no changes to memory",
			},
		},
		memSize:       16,
		frameSize:     8,
		freeListState: []bool{true, false},
	},
	{
		operations: []TAllocMultipleOperation{
			{
				in: 15, pid: 0,
				wantPageTable: &PageTable{[]int{1, 3}},
				wantFreeList:  []bool{false, false, false, false},
				wantError:     nil,
				desc:          "allocate 2 frames that are not contiguous in memory layout",
			},
		},
		memSize:       32,
		frameSize:     8,
		freeListState: []bool{false, true, false, true},
	},
	{
		operations: []TAllocMultipleOperation{ // Allocate to several processes in sequence
			{
				in: 12, pid: 0,
				wantPageTable: &PageTable{[]int{0, 1}},
				wantFreeList:  []bool{false, false, true, true, true, true, true, true},
				wantError:     nil,
				desc:          "Allocate 2 frames to process 0",
			},
			{
				in: 8, pid: 1,
				wantPageTable: &PageTable{[]int{2}},
				wantFreeList:  []bool{false, false, false, true, true, true, true, true},
				wantError:     nil,
				desc:          "Allocate 1 frame to process 1",
			},
			{
				in: 1, pid: 2,
				wantPageTable: &PageTable{[]int{3}},
				wantFreeList:  []bool{false, false, false, false, true, true, true, true},
				wantError:     nil,
				desc:          "Allocate 1 frame to process 2",
			},
			{
				in: 8, pid: 0,
				wantPageTable: &PageTable{[]int{0, 1, 4}},
				wantFreeList:  []bool{false, false, false, false, false, true, true, true},
				wantError:     nil,
				desc:          "Allocate an additional frame to process 0. The frame is not contiguous in memory to the rest of process 0's address space.",
			},
			{
				in: 32, pid: 1,
				wantPageTable: &PageTable{[]int{2}},
				wantFreeList:  []bool{false, false, false, false, false, true, true, true},
				wantError:     errAllocNotEnoughFrames,
				desc:          "Process 1 tries to allocate 32 bytes (4 frames) when only 24 bytes (3 frames) are available -> error",
			},
			{
				in: 16, pid: 1,
				wantPageTable: &PageTable{[]int{2, 5, 6}},
				wantFreeList:  []bool{false, false, false, false, false, false, false, true},
				wantError:     nil,
				desc:          "Allocate an additional frame to process 1",
			},
			{
				in: 8, pid: 2,
				wantPageTable: &PageTable{[]int{3, 7}},
				wantFreeList:  []bool{false, false, false, false, false, false, false, false},
				wantError:     nil,
				desc:          "Allocate the final frame to process 2",
			},
			{
				in: 2, pid: 0,
				wantPageTable: &PageTable{[]int{0, 1, 4}},
				wantFreeList:  []bool{false, false, false, false, false, false, false, false},
				wantError:     errAllocNotEnoughFrames,
				desc:          "Process 1 tries to allocate 2 bytes (rounded up to 1 frame) when no more memory is available -> error",
			},
		},
		memSize:       64,
		frameSize:     8,
		freeListState: []bool{true, true, true, true, true, true, true, true},
	},
}

var ReadTests = []struct {
	addr               int        // virtual address to start reading from
	n                  int        // number of bytes to read
	pageTable          *PageTable // state of process' page table
	frames             [][]byte   // state of MMU.frames
	memSize, frameSize int        // parameters to NewMMU
	content            []byte     // expected memory content output
	err                error      // expected error output
	desc               string     // description of this test
}{
	{
		addr: 0, n: 1,
		pageTable: nil,
		frames:    [][]byte{{0}},
		memSize:   1, frameSize: 1,
		content: nil,
		err:     errReadInvalidAddr,
		desc:    "process tries to read address 0, but has no allocated memory -> error",
	},
	{
		addr: 0, n: 1,
		pageTable: &PageTable{},
		frames:    [][]byte{{0}},
		memSize:   1, frameSize: 1,
		content: nil,
		err:     errReadInvalidAddr,
		desc:    "process tries to read address 0, but has no allocated memory -> error",
	},
	{
		addr: 0x4, n: 1,
		pageTable: &PageTable{[]int{0}},
		frames: [][]byte{
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		memSize: 8, frameSize: 4,
		content: nil,
		err:     errReadInvalidAddr,
		desc:    "process tries to read virtual address 4 (-> page 1, offset 0), but only has virtual address space from 0 to 3 -> error",
	},
	{
		addr: 0x3, n: 2,
		pageTable: &PageTable{[]int{0}},
		frames: [][]byte{
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		memSize: 8, frameSize: 4,
		content: nil,
		err:     errReadOutOfBounds,
		desc:    "process tries to read 2 bytes from virtual address 3 (-> page 0, offset 3), but only has virtual address space from 0 to 3 -> out of bounds error",
	},
	{
		addr: 0x0, n: 5,
		pageTable: &PageTable{[]int{0}},
		frames: [][]byte{
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		memSize: 8, frameSize: 4,
		content: nil,
		err:     errReadOutOfBounds,
		desc:    "process tries to read 5 bytes from virtual address 0, but only has virtual address space from 0 to 3 -> out of bounds error",
	},
	{
		addr: 0x0, n: 1,
		pageTable: &PageTable{[]int{0}},
		frames: [][]byte{
			{1, 0, 0, 0},
			{0, 0, 0, 0},
		},
		memSize: 8, frameSize: 4,
		content: []byte{1},
		err:     nil,
		desc:    "process reads 1 byte from virtual address 0",
	},
	{
		addr: 0x0, n: 1,
		pageTable: &PageTable{[]int{1}},
		frames: [][]byte{
			{1, 0, 0, 0},
			{2, 0, 0, 0},
		},
		memSize: 8, frameSize: 4,
		content: []byte{2},
		err:     nil,
		desc:    "process reads 1 byte from virtual address 0 (-> frame 1, offset 0)",
	},
	{
		addr: 0x0, n: 2,
		pageTable: &PageTable{[]int{1}},
		frames: [][]byte{
			{0, 0, 0, 0},
			{1, 2, 0, 0},
		},
		memSize: 8, frameSize: 4,
		content: []byte{1, 2},
		err:     nil,
		desc:    "process reads 2 bytes starting from virtual address 0 (-> frame 1, offset 0)",
	},
	{
		addr: 0x0, n: 8,
		pageTable: &PageTable{[]int{0, 1}},
		frames: [][]byte{
			{1, 0, 0, 2},
			{3, 0, 0, 4},
		},
		memSize: 8, frameSize: 4,
		content: []byte{1, 0, 0, 2, 3, 0, 0, 4},
		err:     nil,
		desc:    "process reads 8 bytes starting from virtual address 0",
	},
	{
		addr: 0x3, n: 2,
		pageTable: &PageTable{[]int{0, 1}},
		frames: [][]byte{
			{1, 0, 0, 2},
			{3, 0, 0, 4},
		},
		memSize: 8, frameSize: 4,
		content: []byte{2, 3},
		err:     nil,
		desc:    "process reads 2 bytes starting from virtual address 3 (-> {page 0, offset 3} to {page 1, offset 0})",
	},
	{
		addr: 0x3, n: 2,
		pageTable: &PageTable{[]int{1, 0}},
		frames: [][]byte{
			{1, 0, 0, 2},
			{3, 0, 0, 4},
		},
		memSize: 8, frameSize: 4,
		content: []byte{4, 1},
		err:     nil,
		desc:    "process reads 2 bytes starting from virtual address 3 (-> {page 0, offset 3} to {page 1, offset 0})",
	},
	{
		addr: 0x0, n: 4,
		pageTable: &PageTable{[]int{1, 7, 3, 5}},
		frames: [][]byte{
			{0},
			{1},
			{2},
			{3},
			{4},
			{7},
			{6},
			{3},
		},
		memSize: 8, frameSize: 1,
		content: []byte{1, 3, 3, 7},
		err:     nil,
		desc:    "process reads 4 bytes starting from virtual address 0, where frame size = 1; reads 4 bytes spread across physical memory",
	},
	{
		addr: 0x7, n: 5,
		pageTable: &PageTable{[]int{2, 7, 0, 3, 4, 6}},
		frames: [][]byte{
			{0, 0},
			{0, 0},
			{0, 0},
			{0, 1},
			{2, 3},
			{0, 0},
			{4, 5},
			{0, 0},
		},
		memSize: 8, frameSize: 2,
		content: []byte{1, 2, 3, 4, 5},
		err:     nil,
		desc:    "process reads 5 bytes starting from virtual address 7 (-> page 3, offset 1), where frame size = 2; reads 5 bytes spread across 3 frames in physical memory",
	},
}

var WriteTests = []struct {
	name               string
	content            []byte     // content to write
	addr               int        // virtual address to write message
	pageTable          *PageTable // state of process' page table before writing
	frames             [][]byte   // state of MMU.frames before writing
	freeList           []bool     // state of MMU.freeList before writing
	memSize, frameSize int        // parameters to NewMMU
	err                error      // expected error output
	wantPageTable      *PageTable // state of process' page table after writing
	wantFrames         [][]byte   // state of MMU.frames after writing
	wantFreeList       []bool     // state of MMU.freeList after writing
	desc               string     // description of this test
}{
	{
		name:          "no_alloc_write",
		content:       []byte{0},
		addr:          0x00,
		pageTable:     nil,
		frames:        [][]byte{{0}},
		freeList:      []bool{true},
		err:           errWriteOutOfBounds,
		wantPageTable: nil,
		wantFrames:    [][]byte{{0}},
		wantFreeList:  []bool{true},
		memSize:       1, frameSize: 1,
		desc: "process has no memory allocated and tries to write -> error",
	},
	{
		name:          "write_byte",
		content:       []byte{1},
		addr:          0x00,
		pageTable:     &PageTable{[]int{0}},
		frames:        [][]byte{{0}},
		freeList:      []bool{false},
		err:           nil,
		wantPageTable: &PageTable{[]int{0}},
		wantFrames:    [][]byte{{1}},
		wantFreeList:  []bool{false},
		memSize:       1, frameSize: 1,
		desc: "process writes 1 to virtual address 0 (0x00) -> memory updated",
	},
	{
		name:          "write_2_bytes",
		content:       []byte{1, 2},
		addr:          0x00,
		pageTable:     &PageTable{[]int{0, 1}},
		frames:        [][]byte{{0}, {0}},
		freeList:      []bool{false, false},
		err:           nil,
		wantPageTable: &PageTable{[]int{0, 1}},
		wantFrames:    [][]byte{{1}, {2}},
		wantFreeList:  []bool{false, false},
		memSize:       2, frameSize: 1,
		desc: "process writes 1 and 2 starting at virtual address 0 (0x00) -> memory updated at virtual address 0 and 1, which point to physical addresses 0 and 1",
	},
	{
		name:          "write_2_bytes_alt",
		content:       []byte{1, 2},
		addr:          0x00,
		pageTable:     &PageTable{[]int{1, 0}},
		frames:        [][]byte{{0}, {0}},
		freeList:      []bool{false, false},
		err:           nil,
		wantPageTable: &PageTable{[]int{1, 0}},
		wantFrames:    [][]byte{{2}, {1}},
		wantFreeList:  []bool{false, false},
		memSize:       2, frameSize: 1,
		desc: "process writes 1 and 2 starting at virtual address 0 (0x00) \n-> memory updated at virtual address 0 and 1, which point to physical addresses 1 and 0",
	},
	{
		name:      "write_byte_larger_mem",
		content:   []byte{1},
		addr:      0x00,
		pageTable: &PageTable{[]int{0}},
		frames: [][]byte{
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		freeList:      []bool{false, true},
		err:           nil,
		wantPageTable: &PageTable{[]int{0}},
		wantFrames: [][]byte{
			{1, 0, 0, 0},
			{0, 0, 0, 0},
		},
		wantFreeList: []bool{false, true},
		memSize:      8, frameSize: 4,
		desc: "process writes 1 starting at virtual address 0 (0x00) -> memory updated at virtual address 0 (points to frames[0][0])",
	},
	{
		name:      "write_byte_larger_mem_with_offset",
		content:   []byte{1},
		addr:      0x4,
		pageTable: &PageTable{[]int{0, 1}},
		frames: [][]byte{
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		freeList:      []bool{false, false},
		err:           nil,
		wantPageTable: &PageTable{[]int{0, 1}},
		wantFrames: [][]byte{
			{0, 0, 0, 0},
			{1, 0, 0, 0},
		},
		wantFreeList: []bool{false, false},
		memSize:      8, frameSize: 4,
		desc: "process writes 1 starting at virtual address 4 (0x4), frame size = 4 -> memory updated at virtual address 4 (points to frames[1][0])",
	},
	{
		name:      "write_byte_middle_of_frame",
		content:   []byte{1},
		addr:      0x5,
		pageTable: &PageTable{[]int{0, 1}},
		frames: [][]byte{
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		freeList:      []bool{false, false},
		err:           nil,
		wantPageTable: &PageTable{[]int{0, 1}},
		wantFrames: [][]byte{
			{0, 0, 0, 0},
			{0, 1, 0, 0},
		},
		wantFreeList: []bool{false, false},
		memSize:      8, frameSize: 4,
		desc: "process writes 1 starting at virtual address 5 (0x5) -> memory updated at virtual address 5 (points to frames[1][1])",
	},
	{
		name:      "write_byte_offset_frame",
		content:   []byte{1},
		addr:      0x6,
		pageTable: &PageTable{[]int{1, 0}},
		frames: [][]byte{
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		freeList:      []bool{false, false},
		err:           nil,
		wantPageTable: &PageTable{[]int{1, 0}},
		wantFrames: [][]byte{
			{0, 0, 1, 0},
			{0, 0, 0, 0},
		},
		wantFreeList: []bool{false, false},
		memSize:      8, frameSize: 4,
		desc: "process writes 1 starting at virtual address 6 (0x6) -> memory updated at virtual address 5 (points to physical address 2 (frames[0][2]))",
	},
	{
		name:      "write_byte_offset",
		content:   []byte{1},
		addr:      0x0,
		pageTable: &PageTable{[]int{1}},
		frames: [][]byte{
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		freeList:      []bool{true, false},
		err:           nil,
		wantPageTable: &PageTable{[]int{1}},
		wantFrames: [][]byte{
			{0, 0, 0, 0},
			{1, 0, 0, 0},
		},
		wantFreeList: []bool{true, false},
		memSize:      8, frameSize: 4,
		desc: "process writes 1 starting at virtual address 0 (0x0) -> memory updated at virtual address 0 (points to physical address 4 (frames[1][0]))",
	},
	{
		name:      "write_5_bytes_OOM",
		content:   []byte{1, 2, 3, 4, 5},
		addr:      0x0,
		pageTable: &PageTable{[]int{1}},
		frames: [][]byte{
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		freeList:      []bool{false, false},
		err:           errAllocNotEnoughFrames,
		wantPageTable: &PageTable{[]int{1}},
		wantFrames: [][]byte{
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		wantFreeList: []bool{false, false},
		memSize:      8, frameSize: 4,
		desc: "process tries to write '12345' starting at virtual address 0 (0x0), but cannot be allocated enough frames -> error and no changes to memory",
	},
	{
		name:      "write_5_bytes",
		content:   []byte{1, 2, 3, 4, 5},
		addr:      0x0,
		pageTable: &PageTable{[]int{1}},
		frames: [][]byte{
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		freeList:      []bool{true, false},
		err:           nil,
		wantPageTable: &PageTable{[]int{1, 0}},
		wantFrames: [][]byte{
			{5, 0, 0, 0},
			{1, 2, 3, 4},
		},
		wantFreeList: []bool{false, false},
		memSize:      8, frameSize: 4,
		desc: "process tries to write '12345' starting at virtual address 0 (0x0) -> is allocated 1 frame, writes to frames[1][...] and then frames[0][0] (virtual address 0 points to physical frame 1)",
	},
	{
		name:      "write_8_bytes",
		content:   []byte{1, 2, 3, 4, 5, 6, 7, 8},
		addr:      0x5,
		pageTable: &PageTable{[]int{4, 5, 2}},
		frames: [][]byte{
			{0, 0},
			{0, 0},
			{0, 0},
			{0, 0},
			{0, 0},
			{0, 0},
			{0, 0},
			{0, 0},
		},
		freeList:      []bool{true, false, false, true, false, false, true, true},
		err:           nil,
		wantPageTable: &PageTable{[]int{4, 5, 2, 0, 3, 6, 7}},
		wantFrames: [][]byte{
			{2, 3},
			{0, 0},
			{0, 1},
			{4, 5},
			{0, 0},
			{0, 0},
			{6, 7},
			{8, 0},
		},
		wantFreeList: []bool{false, false, false, false, false, false, false, false},
		memSize:      16, frameSize: 2,
		desc: `Process writes '12345678' starting at address 5 (-> virtual page 2, offset 1) -> is allocated 4 frames spread in memory, successful write
                        addr: 0x5 = 0101
			frameSize = 2^1 -> 010 | 1 -> page 2, offset 1
			pagetable and contents (before):
			0 -> 4 (0 0)
			1 -> 5 (0 0)
			2 -> 2 (0 0)
			          ^starting position
			content: (1 2 3 4 5 6 7 8)
			need alloc: 7 bytes -> 4 frames
			free list: (1 0 0 1 0 0 1 1)
			frames to alloc to process: (0 3 6 7)
			new pagetable and contents:
			0 -> 4 (0 0)
			1 -> 5 (0 0)
			2 -> 2 (0 1)
			3 -> 0 (2 3)
			4 -> 3 (4 5)
			5 -> 6 (6 7)
			6 -> 7 (8 0)`,
	},
}

var FreeTests = []struct {
	// input
	pid, n             int
	memSize, frameSize int // parameters to NewMMU
	// initial state
	processes map[int]*PageTable
	freeList  []bool
	frames    [][]byte
	// results
	err           error
	wantProcesses map[int]*PageTable
	wantFreeList  []bool
	wantFrames    [][]byte
	// description
	desc string
}{
	{
		pid: 0, n: 1, memSize: 2, frameSize: 1,
		processes:     map[int]*PageTable{},
		freeList:      []bool{true, true},
		frames:        [][]byte{{0}, {0}},
		err:           errFreeTooManyPages,
		wantProcesses: map[int]*PageTable{},
		wantFreeList:  []bool{true, true},
		wantFrames:    [][]byte{{0}, {0}},
		desc:          "process 0 tries to free before having allocated any memory -> error",
	},
	{
		pid: 0, n: 1, memSize: 2, frameSize: 1,
		processes:     map[int]*PageTable{0: {[]int{}}},
		freeList:      []bool{true, true},
		frames:        [][]byte{{1}, {0}},
		err:           errFreeTooManyPages,
		wantProcesses: map[int]*PageTable{0: {[]int{}}},
		wantFreeList:  []bool{true, true},
		wantFrames:    [][]byte{{1}, {0}},
		desc:          "process 0 tries to free 1 page, but has 0 allocated -> error",
	},
	{
		pid: 0, n: 2, memSize: 2, frameSize: 1,
		processes:     map[int]*PageTable{0: {[]int{0}}},
		freeList:      []bool{false, true},
		frames:        [][]byte{{1}, {0}},
		err:           errFreeTooManyPages,
		wantProcesses: map[int]*PageTable{0: {[]int{0}}},
		wantFreeList:  []bool{false, true},
		wantFrames:    [][]byte{{1}, {0}},
		desc:          "process 0 tries to free 2 pages, only has 1 allocated -> error",
	},
	{
		pid: 0, n: 1, memSize: 2, frameSize: 1,
		processes:     map[int]*PageTable{0: {[]int{0}}},
		freeList:      []bool{false, true},
		frames:        [][]byte{{1}, {0}},
		err:           nil,
		wantProcesses: map[int]*PageTable{0: {[]int{}}},
		wantFreeList:  []bool{true, true},
		wantFrames:    [][]byte{{0}, {0}},
		desc:          "process 0 frees 1 page (-> frame 0) -> page table and free list updated, free memory set to 0",
	},
	{
		pid: 0, n: 1, memSize: 2, frameSize: 1,
		processes:     map[int]*PageTable{0: {[]int{1}}},
		freeList:      []bool{true, false},
		frames:        [][]byte{{0}, {1}},
		err:           nil,
		wantProcesses: map[int]*PageTable{0: {[]int{}}},
		wantFreeList:  []bool{true, true},
		wantFrames:    [][]byte{{0}, {0}},
		desc:          "process 0 frees 1 page (-> frame 1) -> page table and free list updated, free memory set to 0",
	},
	{
		pid: 0, n: 2, memSize: 2, frameSize: 1,
		processes:     map[int]*PageTable{0: {[]int{0, 1}}},
		freeList:      []bool{false, false},
		frames:        [][]byte{{1}, {2}},
		err:           nil,
		wantProcesses: map[int]*PageTable{0: {[]int{}}},
		wantFreeList:  []bool{true, true},
		wantFrames:    [][]byte{{0}, {0}},
		desc:          "process 0 frees 2 pages (-> frames 0, 1) -> page table and free list updated, free memory set to 0",
	},
	{
		pid: 0, n: 2, memSize: 2, frameSize: 1,
		processes:     map[int]*PageTable{0: {[]int{1, 0}}},
		freeList:      []bool{false, false},
		frames:        [][]byte{{1}, {2}},
		err:           nil,
		wantProcesses: map[int]*PageTable{0: {[]int{}}},
		wantFreeList:  []bool{true, true},
		wantFrames:    [][]byte{{0}, {0}},
		desc:          "process 0 frees 2 pages (-> frames 1, 0) -> page table and free list updated, free memory set to 0",
	},
	{
		pid: 0, n: 4, memSize: 8, frameSize: 1,
		processes:     map[int]*PageTable{0: {[]int{0, 2, 4, 6}}},
		freeList:      []bool{false, true, false, true, false, true, false, true},
		frames:        [][]byte{{1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}},
		err:           nil,
		wantProcesses: map[int]*PageTable{0: {[]int{}}},
		wantFreeList:  []bool{true, true, true, true, true, true, true, true},
		wantFrames:    [][]byte{{0}, {2}, {0}, {4}, {0}, {6}, {0}, {8}},
		desc:          "process 0 frees 4 pages (-> frames 0, 2, 4, 6) -> page table and free list updated, free memory set to 0",
	},
	{
		pid: 1, n: 2, memSize: 16, frameSize: 2,
		processes: map[int]*PageTable{0: {[]int{0, 2}}, 1: {[]int{5, 3, 1}}, 2: {[]int{4, 6}}},
		freeList:  []bool{false, false, false, false, false, false, false, true},
		frames: [][]byte{
			{0, 1},
			{104, 105},
			{2, 3},
			{102, 103},
			{200, 201},
			{100, 101},
			{202, 203},
			{0, 0},
		},
		err:           nil,
		wantProcesses: map[int]*PageTable{0: {[]int{0, 2}}, 1: {[]int{5}}, 2: {[]int{4, 6}}},
		wantFreeList:  []bool{false, true, false, true, false, false, false, true},
		wantFrames: [][]byte{
			{0, 1},
			{0, 0}, // <--- freed
			{2, 3},
			{0, 0}, // <--- freed
			{200, 201},
			{100, 101},
			{202, 203},
			{0, 0},
		},
		desc: "process 1 frees 2 pages (-> frames 3, 1) -> page table and free list updated, free memory set to 0, other processes unaffected",
	},
	{
		pid: 2, n: 3, memSize: 32, frameSize: 4,
		processes: map[int]*PageTable{0: {[]int{0, 2}}, 1: {[]int{5, 3}}, 2: {[]int{4, 1, 6, 7}}},
		freeList:  []bool{false, false, false, false, false, false, false, false},
		frames: [][]byte{
			{0, 1, 2, 3},
			{204, 205, 206, 207},
			{4, 5, 6, 7},
			{104, 105, 106, 107},
			{200, 201, 202, 203},
			{100, 101, 102, 103},
			{208, 209, 210, 211},
			{212, 213, 214, 215},
		},
		err:           nil,
		wantProcesses: map[int]*PageTable{0: {[]int{0, 2}}, 1: {[]int{5, 3}}, 2: {[]int{4}}},
		wantFreeList:  []bool{false, true, false, false, false, false, true, true},
		wantFrames: [][]byte{
			{0, 1, 2, 3},
			{0, 0, 0, 0}, // <--- freed
			{4, 5, 6, 7},
			{104, 105, 106, 107},
			{200, 201, 202, 203},
			{100, 101, 102, 103},
			{0, 0, 0, 0}, // <--- freed
			{0, 0, 0, 0}, // <--- freed
		},
		desc: "process 2 frees 3 pages (-> frames 1, 6, 7) -> page table and free list updated, free memory set to 0, other processes unaffected",
	},
}

// command types for 'TestSequences'
const (
	cmdAlloc = iota
	cmdWrite
	cmdRead
	cmdFree
)

type mmuCmd struct {
	// command to perform; Alloc, Write, Read or Free
	cmd int
	// process id
	pid int
	n   int
	// virtualAddr: input to mmu.Write and mmu.Read
	virtualAddr int
	// input to mmu.Write or output from mmu.Read
	content []byte
	// resulting error
	err error
	// description of this step
	desc string
}

func (cmd mmuCmd) String() string {
	var str string
	switch cmd.cmd {
	case cmdAlloc:
		str = fmt.Sprintf("Alloc(pid = %d, n = %d)", cmd.pid, cmd.n)
	case cmdWrite:
		str = fmt.Sprintf("Write(pid = %d, virtualAddress = %d, content = %v)", cmd.pid, cmd.virtualAddr, cmd.content)
	case cmdRead:
		str = fmt.Sprintf("Read(pid = %d, virtualAddress = %d, n = %d)", cmd.pid, cmd.virtualAddr, cmd.n)
	case cmdFree:
		str = fmt.Sprintf("Free(pid = %d, n = %d)", cmd.pid, cmd.n)
	}
	return str
}

var SequenceTests = []struct {
	memSize, frameSize int      // parameters to NewMMU
	cmds               []mmuCmd // each command is performed in sequence
	// expected final state after all commands are complete
	wantFrames    [][]byte
	wantFreeList  []bool
	wantProcesses map[int]*PageTable
}{
	{
		// tests:
		// - allocating and writing to memory, starting in the middle of a frame
		// - freeing parts of memory and having it allocated by another process
		// - checks that freed memory content has been set to 0
		memSize: 16, frameSize: 8,
		cmds: []mmuCmd{
			{cmd: cmdAlloc, pid: 1, n: 13, desc: "process 1 allocates 13 bytes (-> 16 bytes, 2 frames)"},
			{cmd: cmdWrite, pid: 1, virtualAddr: 6, content: []byte{1, 2, 3, 4, 5}, desc: "process 1 writes {1 2 3 4 5} starting at virtual address 6"},
			{cmd: cmdFree, pid: 1, n: 1, desc: "process 1 frees 1 page (frame #1)"},
			{cmd: cmdAlloc, pid: 2, n: 8, desc: "process 2 allocates 8 bytes (-> 1 frame), which will be frame #1"},
			{cmd: cmdRead, pid: 2, virtualAddr: 0, n: 8, content: []byte{0, 0, 0, 0, 0, 0, 0, 0}, desc: "process 2 reads 8 bytes from which were previously allocated to process 1; content is 0 as expected after process 1 freed it"},
		},
		wantFrames: [][]byte{
			{0, 0, 0, 0, 0, 0, 1, 2},
			{0, 0, 0, 0, 0, 0, 0, 0},
		},
		wantFreeList: []bool{false, false},
		wantProcesses: map[int]*PageTable{
			1: {[]int{0}},
			2: {[]int{1}},
		},
	},
	{
		// tests:
		// - allocating memory and then freeing parts of memory
		// - allocating fragmented frames
		memSize: 32, frameSize: 4,
		cmds: []mmuCmd{
			{cmd: cmdAlloc, pid: 1, n: 12, desc: "process 1 allocates 12 bytes (-> 3 frames)"},
			{cmd: cmdAlloc, pid: 2, n: 12, desc: "process 2 allocates 13 bytes (-> 3 frames)"},
			{cmd: cmdFree, pid: 1, n: 1, desc: "process 1 frees 1 page (page 2 -> frame 2)"},
			{cmd: cmdFree, pid: 2, n: 1, desc: "process 2 frees 1 page (page 2 -> frame 5)"},
			{cmd: cmdAlloc, pid: 3, n: 16, desc: "process 3 allocates 16 bytes (-> 4 frames, allocated frames 2, 5, 6, 7)"},
		},
		wantFrames: [][]byte{
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
		wantFreeList: []bool{false, false, false, false, false, false, false, false},
		wantProcesses: map[int]*PageTable{
			1: {[]int{0, 1}},
			2: {[]int{3, 4}},
			3: {[]int{2, 5, 6, 7}},
		},
	},
	{
		// tests:
		// - allocating 1 frame, then writing 3 frames worth of content -> dynamically allocating 2 more frames
		// - another process tries to allocate 7 frames -> failure since only 5 free frames left due to dynamic allocation
		// - the other process instead allocates 5 frames and writes content to the middle of its address space
		memSize: 32, frameSize: 4,
		cmds: []mmuCmd{
			{cmd: cmdAlloc, pid: 1, n: 4, desc: "process 1 allocates 4 bytes (-> 1 frame)"},
			{cmd: cmdWrite, pid: 1, virtualAddr: 0, content: []byte("0123456789ab"), desc: "process 1 writes 12 bytes (-> 3 frames), is dynamically allocated 8 bytes"},
			{cmd: cmdAlloc, pid: 2, n: 28, err: errAllocNotEnoughFrames, desc: "process 2 tries to allocate 28 bytes (-> 7 frames), which would be free if process 1 did not dynamically allocate during previous Write -> error"},
			{cmd: cmdAlloc, pid: 2, n: 20, desc: "process 2 allocates the remaining free 20 bytes (-> 5 frames)"},
			{cmd: cmdWrite, pid: 2, virtualAddr: 8, content: []byte{100, 101, 102, 103, 104, 105, 106, 107}, desc: "process 2 writes 8 bytes starting at virtual address 8 (-> page 2, offset 0)"},
		},
		wantFrames: [][]byte{
			[]byte("0123"),
			[]byte("4567"),
			[]byte("89ab"),
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{100, 101, 102, 103},
			{104, 105, 106, 107},
			{0, 0, 0, 0},
		},
		wantFreeList: []bool{false, false, false, false, false, false, false, false},
		wantProcesses: map[int]*PageTable{
			1: {[]int{0, 1, 2}},
			2: {[]int{3, 4, 5, 6, 7}},
		},
	},
	{
		// tests: allocating, freeing, writing and reading with 5 processes and fragmented memory
		memSize: 32, frameSize: 2,
		cmds: []mmuCmd{
			{cmd: cmdAlloc, pid: 1, n: 4, desc: "process 1 allocates 4 bytes (-> 2 frames; frames 0, 1)"},
			{cmd: cmdAlloc, pid: 2, n: 8, desc: "process 2 allocates 8 bytes (-> 4 frames; frames 2, 3, 4, 5)"},
			{cmd: cmdAlloc, pid: 3, n: 4, desc: "process 3 allocates 4 bytes (-> 2 frames; frames 6, 7)"},
			{cmd: cmdAlloc, pid: 4, n: 8, desc: "process 4 allocates 8 bytes (-> 4 frames; frames 8, 9, 10, 11)"},
			{cmd: cmdAlloc, pid: 5, n: 8, desc: "process 5 allocates 8 bytes (-> 4 frames; frames 12, 13, 14, 15)"},
			{cmd: cmdWrite, pid: 3, virtualAddr: 1, content: []byte{30, 31, 32}, desc: "process 3 writes 3 bytes starting from virtual address 1"},
			{cmd: cmdFree, pid: 2, n: 3, desc: "process 2 frees 3 pages (-> frames 3, 4, 5)"},
			// after above operation: frames 3, 4, 5 free
			{cmd: cmdWrite, pid: 1, virtualAddr: 3, content: []byte{10, 11, 12, 13, 14}, desc: "process 1 writes 5 bytes starting from virtual address 3 (-> page 2, offset 1) -> dynamically allocated 2 frames (-> frames 4, 5)"},
			// after above operation: frame 5 free
			{cmd: cmdRead, pid: 1, virtualAddr: 0, n: 8, content: []byte{0, 0, 0, 10, 11, 12, 13, 14}, desc: "process 1 reads all of its memory content"},
			{cmd: cmdFree, pid: 4, n: 3, desc: "process 4 frees 3 pages (-> frames 9, 10, 11)"},
			// after above operation: frames 5, 9, 10, 11 free
			{cmd: cmdWrite, pid: 2, virtualAddr: 0, content: []byte{20, 21, 22, 23, 24, 25}, desc: "process 2 writes 6 bytes starting at virtual address 0 -> is dynamically allocated frames 5, 9; writes (in this order) to physical frames 2, 5, 9"},
			// after above operation: frames 10, 11 free
			{cmd: cmdWrite, pid: 5, virtualAddr: 0, content: []byte{50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61}, desc: "process 5 writes 12 bytes starting at virtual address 0 -> is dynamically allocated frames 10, 11; writes (in this order) to physical frames 12, 13, 14, 15, 10, 11"},
			// after above operation: no free frames
			{cmd: cmdRead, pid: 1, virtualAddr: 0, n: 8, content: []byte{0, 0, 0, 10, 11, 12, 13, 14}, desc: "process 1 reads all of its memory content"},
			{cmd: cmdRead, pid: 2, virtualAddr: 0, n: 6, content: []byte{20, 21, 22, 23, 24, 25}, desc: "process 2 reads all of its memory content"},
			{cmd: cmdRead, pid: 3, virtualAddr: 0, n: 4, content: []byte{0, 30, 31, 32}, desc: "process 3 reads all of its memory content"},
			{cmd: cmdRead, pid: 5, virtualAddr: 0, n: 12, content: []byte{50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61}, desc: "process 5 reads all of its memory content"},
		},
		wantFrames: [][]byte{
			{0, 0},   // process 1
			{0, 10},  // process 1
			{20, 21}, // process 2
			{11, 12}, // process 1
			{13, 14}, // process 1
			{22, 23}, // process 2
			{0, 30},  // process 3
			{31, 32}, // process 3
			{0, 0},   // process 4
			{24, 25}, // process 2
			{58, 59}, // process 5
			{60, 61}, // process 5
			{50, 51}, // process 5
			{52, 53}, // process 5
			{54, 55}, // process 5
			{56, 57}, // process 5
		},
		wantFreeList: []bool{false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false},
		wantProcesses: map[int]*PageTable{
			1: {[]int{0, 1, 3, 4}},
			2: {[]int{2, 5, 9}},
			3: {[]int{6, 7}},
			4: {[]int{8}},
			5: {[]int{12, 13, 14, 15, 10, 11}},
		},
	},
}

// --------------------
// | page table tests |
// --------------------

var (
	errPTFree   = errors.New("tried to free unallocated pages")
	errPTLookup = errors.New("tried to read out of bounds page")
)

type TPTFreeWant struct {
	pageTable []int
	freed     []int
	err       error
}

type TPTLookupWant struct {
	frameIndex int
	err        error
}

var PTAppendTests = []struct {
	in        []int
	want      []int
	pageTable []int
}{
	{
		in:        []int{},
		want:      []int{},
		pageTable: []int{},
	},
	{
		in:        []int{0},
		want:      []int{0},
		pageTable: []int{},
	},
	{
		in:        []int{0},
		want:      []int{1, 0},
		pageTable: []int{1},
	},
	{
		in:        []int{0, 2, 4},
		want:      []int{1, 3, 0, 2, 4},
		pageTable: []int{1, 3},
	},
}

var PTFreeTests = []struct {
	pageTable []int
	in        int
	want      TPTFreeWant
}{
	{
		pageTable: []int{},
		in:        1,
		want: TPTFreeWant{
			pageTable: []int{},
			freed:     []int{},
			err:       errPTFree,
		},
	},
	{
		pageTable: []int{0},
		in:        1,
		want: TPTFreeWant{
			pageTable: []int{},
			freed:     []int{0},
			err:       nil,
		},
	},
	{
		pageTable: []int{0, 1},
		in:        1,
		want: TPTFreeWant{
			pageTable: []int{0},
			freed:     []int{1},
			err:       nil,
		},
	},
	{
		pageTable: []int{0, 1},
		in:        2,
		want: TPTFreeWant{
			pageTable: []int{},
			freed:     []int{0, 1},
			err:       nil,
		},
	},
	{
		pageTable: []int{3, 2, 1},
		in:        2,
		want: TPTFreeWant{
			pageTable: []int{3},
			freed:     []int{2, 1},
			err:       nil,
		},
	},
	{
		pageTable: []int{3, 2, 1},
		in:        4,
		want: TPTFreeWant{
			pageTable: []int{3, 2, 1},
			freed:     []int{},
			err:       errPTFree,
		},
	},
}

var PTLookupTests = []struct {
	pageTable []int
	in        int
	want      TPTLookupWant
}{
	{
		pageTable: []int{},
		in:        0,
		want:      TPTLookupWant{NoEntry, errPTLookup},
	},
	{
		pageTable: []int{1},
		in:        1,
		want:      TPTLookupWant{NoEntry, errPTLookup},
	},
	{
		pageTable: []int{1},
		in:        0,
		want:      TPTLookupWant{1, nil},
	},
	{
		pageTable: []int{3, 2, 1},
		in:        2,
		want:      TPTLookupWant{1, nil},
	},
	{
		pageTable: []int{3, 2, 1},
		in:        0,
		want:      TPTLookupWant{3, nil},
	},
	{
		pageTable: []int{3, 2, 1, 0, 7, 6, 5, 4},
		in:        5,
		want:      TPTLookupWant{6, nil},
	},
}
