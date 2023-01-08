package rr

import (
	"dat320/lab4/scheduler/cpu"
	"dat320/lab4/scheduler/job"
	"time"
)

type roundRobin struct {
	queue   job.Jobs
	cpu     *cpu.CPU
	quantum time.Duration
	//job     *job.Job
}

func New(cpus []*cpu.CPU, quantum time.Duration) *roundRobin {
	// TODO(student) construct new RR scheduler
	if len(cpus) != 1 {
		panic("fifo scheduler supports only a single CPU")
	}
	return &roundRobin{
		cpu:     cpus[0],
		queue:   make(job.Jobs, 0),
		quantum: quantum,
		//job:     cpus[0].CurrentJob(),
	}
}
func (rr *roundRobin) Add(job *job.Job) {
	// TODO(student) Add job to queue
	rr.queue = append(rr.queue, job)
}

// Tick runs the scheduled jobs for the system time, and returns
// the number of jobs finished in this tick. Depending on scheduler requirements,
// the Tick method may assign new jobs to the CPU before returning.
func (rr *roundRobin) Tick(systemTime time.Duration) int {
	jobsFinished := 0
	if rr.cpu.IsRunning() { //cpu is running evaluates to true
		if rr.cpu.Tick() { //cpu.Tick() //increase jobsfinished
			jobsFinished++ //increment the job finished
		}
	}
	if systemTime%rr.quantum == 0 { //time slice exhausted?
		//check to see if current job is finished
		if rr.cpu.CurrentJob() == nil { //check if current job finished
			//jobsFinished++
		} else { //no
			rr.Add(rr.cpu.CurrentJob()) //add to the list of jobs to be run
		}
		rr.reassign() //then reassign
	}
	return jobsFinished
}

// reassign assigns a job to the cpu
func (rr *roundRobin) reassign() {
	nxtJob := rr.getNewJob()
	rr.cpu.Assign(nxtJob)
	// TODO(student) Implement reassign and use it from Tick
}

// getNewJob finds a new job to run on the CPU, removes the job from the queue and returns the job
func (rr *roundRobin) getNewJob() *job.Job {
	if len(rr.queue) == 0 {
		return nil
	}
	removedJob := rr.queue[0]
	rr.queue = rr.queue[1:]
	return removedJob
}
