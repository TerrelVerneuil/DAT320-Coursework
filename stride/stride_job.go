package stride

import (
	"dat320/lab4/scheduler/job"
	"time"
)

// NewJob creates a job for stride scheduling.
func NewJob(size, tickets int, estimated time.Duration) *job.Job {
	const numerator = 10_000
	job := job.New(size, estimated) //creates new job with size and estimated duration
	if tickets > 0 {                //no tickets
		job.Tickets = tickets //assign tickets to job

	}
	job.Stride = (int(estimated) / int(time.Nanosecond*1000)) / 100 // int(estimated)
	// TODO(student) return the job with the correct fields

	return job
}
