package paging

import "errors"

var (
	errOutOfMemory         = errors.New("out of memory")
	errNothingToAllocate   = errors.New("invalid argument: allocation/deallocation request must be greater than 0")
	errNothingToRead       = errors.New("invalid argument: bytes in read request must be greater than 0")
	errAddressOutOfBounds  = errors.New("address out of bounds")
	errIndexOutOfBounds    = errors.New("index out of bounds")
	errFreeListDuplicateOp = errors.New("tried to update a free list entry to its current state")
	errInvalidProcess      = errors.New("process does not exist")
	errFreeOutOfBounds     = errors.New("cannot free more pages than have been allocated")
)

var errNotImplemented = errors.New("this is not yet implemented")
