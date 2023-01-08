package stride

import (
	"dat320/lab4/scheduler/job"
	"dat320/lab4/scheduler/system"
	"testing"
	"time"
)

type minPassTest struct {
	name         string
	jobs         system.Schedule
	passValues   []int
	strideValues []int
	want         int
}

func (tc *minPassTest) setupTest() {
	theJobs := make(system.Schedule, len(tc.jobs))
	copy(theJobs, tc.jobs)
	for i := range tc.jobs {
		theJobs[i].Job.Pass = tc.passValues[i]
		theJobs[i].Job.Stride = tc.strideValues[i]
	}
	tc.jobs = theJobs
}

var passValues = [][]int{
	{0, 0, 0},      // 2
	{0, 0, 40},     // 0
	{100, 0, 40},   // 1
	{100, 200, 40}, // 2
	{100, 200, 80}, // 2
}

var strideValues = []int{100, 200, 40}

const (
	arrival = 0
	t020    = 20 * time.Millisecond
)

var tickets = func(tickets int) *system.Entry {
	return &system.Entry{Job: NewJob(0, tickets, t020), Arrival: arrival}
}

var testCases = []minPassTest{ // stride order: 100, 200, 40
	{"round 1", system.Schedule{tickets(100), tickets(50), tickets(250)}, passValues[0], strideValues, 2},
	{"round 2", system.Schedule{tickets(100), tickets(50), tickets(250)}, passValues[1], strideValues, 0},
	{"round 3", system.Schedule{tickets(100), tickets(50), tickets(250)}, passValues[2], strideValues, 1},
	{"round 4", system.Schedule{tickets(100), tickets(50), tickets(250)}, passValues[3], strideValues, 2},
	{"round 5", system.Schedule{tickets(100), tickets(50), tickets(250)}, passValues[4], strideValues, 2},
}

func TestMinPass(t *testing.T) {
	j := NewJob(0, 50, t020)
	if j == nil {
		t.Fatal("stride.NewJob not implemented")
	}
	for _, tc := range testCases {
		tc.setupTest()
		joblist := make(job.Jobs, len(tc.jobs))
		for i, entry := range tc.jobs {
			joblist[i] = entry.Job
		}
		got := MinPass(joblist)
		if got != tc.want {
			t.Errorf("%s: got job %d, want job %d", tc.name, got, tc.want)
		}
	}
}
