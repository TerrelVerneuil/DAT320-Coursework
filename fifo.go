package fifo

import (
	"dat320/lab4/scheduler/cpu"
	"dat320/lab4/scheduler/job"
	"time"
)

type fifo struct {
	cpu   *cpu.CPU
	queue job.Jobs
}

func New(cpus []*cpu.CPU) *fifo {
	if len(cpus) != 1 {
		panic("fifo scheduler supports only a single CPU")
	}
	return &fifo{
		cpu:   cpus[0],
		queue: make(job.Jobs, 0),
	}
}

func (f *fifo) Add(job *job.Job) {
	f.queue = append(f.queue, job)
}

func (f *fifo) getNewJob() *job.Job {
	if len(f.queue) == 0 {
		return nil
	}
	removedJob := f.queue[0]
	f.queue = f.queue[1:]
	return removedJob
}

// reassign finds a new job to run on this CPU
func (f *fifo) reassign() {
	nxtJob := f.getNewJob()
	f.cpu.Assign(nxtJob)
}

func (f *fifo) Tick(systemTime time.Duration) int {
	jobsFinished := 0
	if f.cpu.IsRunning() {
		if f.cpu.Tick() {
			jobsFinished++
			f.reassign()
		}
	} else {
		//f.Tick(systemTime)
		// CPU is idle, find new job in own queue
		f.reassign()
	}
	return jobsFinished
}
