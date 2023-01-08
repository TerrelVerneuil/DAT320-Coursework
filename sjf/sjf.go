package sjf

import (
	"dat320/lab4/scheduler/cpu"
	"dat320/lab4/scheduler/job"
	"dat320/lab4/scheduler/system/systime"
	"time"
)

type sjf struct {
	queue job.Jobs
	cpu   *cpu.CPU
	//ODO(student) add missing fields, if necessary
	job       *job.Job
	remaining time.Duration
}

func New(cpus []*cpu.CPU) *sjf {
	// ODO(student) construct new RR scheduler

	return &sjf{
		cpu:   cpus[0],
		queue: make(job.Jobs, 0),
		job:   cpus[0].CurrentJob(),
	}
}

func (s *sjf) Add(job *job.Job) {
	// ODO(student) Add job to queue
	s.queue = append(s.queue, job)
}

// Tick runs the scheduled jobs for the system time, and returns
// the number of jobs finished in this tick. Depending on scheduler requirements,
// the Tick method may assign new jobs to the CPU before returning.
func (s *sjf) Tick(systemTime time.Duration) int {
	jobsFinished := 0
	s.remaining -= systemTime * systime.TickDuration
	done := false
	if s.remaining <= 0 {
		done = true
		if done {
			jobsFinished++
			s.reassign()
			return jobsFinished
		}
	} else {
		s.reassign()
	}
	return jobsFinished
}

// reassign assigns a job to the cpu
func (s *sjf) reassign() {
	s.cpu.Assign(s.getNewJob())
	// ODO(student) Implement reassign and use it from Tick
}

// getNewJob finds a new job to run on the CPU, removes the job from the queue and returns the job
func (s *sjf) getNewJob() *job.Job {
	//ODO(student) Implement getNewJob and use it from reassign
	if len(s.queue) == 0 {
		return nil
	}
	removedJob := s.queue[0]
	s.queue = s.queue[1:]
	return removedJob
}
