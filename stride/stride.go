package stride

import (
	"dat320/lab4/scheduler/cpu"
	"dat320/lab4/scheduler/job"
	"dat320/lab4/scheduler/system/systime"
	"time"
)

type stride struct {
	queue job.Jobs
	cpu   *cpu.CPU
	// TODO(student) add missing fields, if necessary
	job       *job.Job
	remaining time.Duration
	quantum   time.Duration
}

func New(cpus []*cpu.CPU, quantum time.Duration) *stride {
	// TODO(student) construct new stride scheduler
	return &stride{
		cpu:     cpus[0],
		quantum: quantum,
		queue:   make(job.Jobs, 0),
		job:     cpus[0].CurrentJob(),
	}
}
func (s *stride) Add(job *job.Job) {
	// TODO(student) Add job to queue
	s.queue = append(s.queue, job)
}

// Tick runs the scheduled jobs for the system time, and returns
// the number of jobs finished in this tick. Depending on scheduler requirements,
// the Tick method may assign new jobs to the CPU before returning.
func (s *stride) Tick(systemTime time.Duration) int {
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
	}
	// TODO(student) Implement Tick()
	return jobsFinished
}

// reassign assigns a job to the cpu
func (s *stride) reassign() {
	// TODO(student) Implement reassign and use it from Tick()
	s.cpu.Assign(s.getNewJob())
}

// getNewJob finds a new job to run on the CPU, removes the job from the queue and returns the job
func (s *stride) getNewJob() *job.Job {
	// TODO(student) Implement getNewJob and use it from reassign
	minimum := MinPass(s.queue) //current                  //pick client with min pass
	s.queue[minimum].Pass += s.queue[minimum].Stride
	return s.queue[minimum]
}

// minPass returns the index of the job with the lowest pass value.
func MinPass(theJobs job.Jobs) int {
	lowest := 0
	for i := 0; i < len(theJobs); i++ {
		if theJobs[lowest].Pass > theJobs[i].Pass { //if the lowest pass is greater than the jobs at i : starts at 0
			lowest = i //assign lowest the lowest value
			continue
		}
		if theJobs[lowest].Pass == theJobs[i].Pass && theJobs[lowest].Stride > theJobs[i].Stride {
			lowest = i
		}
	}

	// TODO(student) Implement MinPass and use it from getNewJob
	return lowest
}
