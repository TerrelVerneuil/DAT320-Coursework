package paging

// Process simulates a (highly simplified) process
type Process struct {
	pid int
	mmu *MMU
}

// NewProcess creates a new process
func NewProcess(pid int, mmu *MMU) *Process {
	return &Process{pid: pid, mmu: mmu}
}

// Malloc requests that the MMU allocates n bytes to this process
func (p *Process) Malloc(n int) (err error) {
	return p.mmu.Alloc(p.pid, n)
}

// Free frees n pages from p, starting from the end of its address space
func (p *Process) Free(n int) {
	if n > 0 {
		_ = p.mmu.Free(p.pid, n)
	}
}

// Read tries to read length bytes starting from virtualAddress
func (p *Process) Read(virtualAddress, length int) (content []byte, err error) {
	return p.mmu.Read(p.pid, virtualAddress, length)
}

// Write tries to write content to the address space of p, starting from virtualAddress
func (p *Process) Write(virtualAddress int, message []byte) (err error) {
	return p.mmu.Write(p.pid, virtualAddress, message)
}
