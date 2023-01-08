package job

import (
	"dat320/lab4/scheduler/system/systime"
	"fmt"
	"time"
)

var nextID = 0

const (
	NotStartedYet   = -1
	DefaultCPUSpeed = 1
)

// Job keeps track of when the job arrived, was started, its remaining time,
// and when it finished. Job also contains the job's working set size,
// its estimated running time and the current speed of the job. The job's
// speed is determined based on whether or not the job is in the cache.
type Job struct {
	id        int
	size      int           // memory size of the job
	estimated time.Duration // the job's duration
	speed     int           // the current speed of the job (cached vs non-cached)
	arrival   time.Duration
	start     time.Duration
	finished  time.Duration
	remaining time.Duration
	systime.SystemTime
	Stride  int
	Pass    int
	Tickets int
}

// New returns a job with given working set size and estimated running time.
func New(size int, estimated time.Duration) *Job {
	nextID++
	return newJob(nextID, size, estimated)
}

// newJob returns a job with given working set size and estimated running time.
func newJob(id, size int, estimated time.Duration) *Job {
	return &Job{
		id:        id,
		size:      size,
		speed:     DefaultCPUSpeed,
		estimated: estimated,
		remaining: estimated,
		start:     NotStartedYet,
	}
}

func NewTestJob(id int, estimated, remaining time.Duration) Job {
	return Job{
		id:        id,
		estimated: estimated,
		remaining: remaining,
	}
}

func (j *Job) Clone() Job {
	if j == nil {
		return *j
	}
	return Job{
		id:        j.id,
		size:      j.size,
		start:     j.start,
		arrival:   j.arrival,
		estimated: j.estimated,
		remaining: j.remaining,
	}
}

func (j Job) Remaining() time.Duration {
	return j.remaining
}

// ID returns the job's ID.
func (j Job) ID() int {
	return j.id
}

// Size returns the size of the job, i.e. how much cache space it takes up.
func (j Job) Size() int {
	return j.size
}

func ResetJobCounter() {
	nextID = 0
}

// SetSpeed sets the job's current speed.
func (j *Job) SetSpeed(speed int) {
	j.speed = speed
}

// run runs the job for the given duration.
func (j *Job) run(durationToRun time.Duration) bool {
	j.remaining -= durationToRun
	return j.remaining <= 0
}

// Tick runs the job for one tick and returns true if job is finished.
func (j *Job) Tick() bool {
	done := j.run(systime.TickDuration * time.Duration(j.speed))
	if done {
		j.finished = j.Now()
	}
	return done
}

// Equal returns true if this job and the given job has the same id.
func (j Job) Equal(job Job) bool {
	return j.id == job.id
}

func (j Job) String() string {
	context := ""
	if j.Pass != 0 || j.Stride != 0 || j.Tickets != 0 {
		context = fmt.Sprintf("(p=%5d,s=%4d,t=%3d)", j.Pass, j.Stride, j.Tickets)
	}
	return fmt.Sprintf("%s(%4v)(%dx) %s", toLetter(j.id), j.remaining, j.speed, context)
}

func (j Job) GoString() string {
	return fmt.Sprintf("jb(%d, %d, %d)", j.id, j.estimated, j.remaining)
}

func (j Job) Name() string {
	return toLetter(j.id)
}
