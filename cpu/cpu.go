package cpu

import (
	"dat320/lab4/scheduler/job"
	"fmt"
)

type CPU struct {
	id      int
	current *job.Job
}

func New(id int) *CPU {
	return &CPU{id: id}
}

func NewCPUs(num int) []*CPU {
	cpus := make([]*CPU, num)
	for i := 0; i < num; i++ {
		cpus[i] = New(i)
	}
	return cpus
}

func (p *CPU) ID() int {
	return p.id
}

func (p *CPU) Assign(job *job.Job) {
	if job != nil {
		job.Started(p.id)
	}
	p.current = job
}

// CurrentJob returns the job currently running on this CPU.
func (p *CPU) CurrentJob() *job.Job {
	return p.current
}

// IsRunning returns true if this CPU is running some job.
// Otherwise, false is returned if the CPU is idle.
func (p *CPU) IsRunning() bool {
	return p.current != nil
}

// Tick runs the current job on this CPU for one clock tick;
// returns true if current job is done.
func (p *CPU) Tick() bool {
	done := p.current.Tick()
	if done {
		// current job is done; mark CPU as idle
		p.current = nil
	}
	return done
}

func (p *CPU) Header() string {
	return fmt.Sprintf("CPU%d", p.ID())
}

func (p *CPU) String() string {
	if p.IsRunning() {
		return p.current.String()
	}
	return "Idle"
}
