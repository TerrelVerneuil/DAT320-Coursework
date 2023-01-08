package job

import (
	"strconv"
	"strings"
)

// Jobs is a slice of jobs ordered according to some scheduling policies.
type Jobs []*Job

// Has returns true if the given job is in the jobs slice.
func (js Jobs) Has(job *Job) bool {
	for _, j := range js {
		if j.Equal(*job) {
			return true
		}
	}
	return false
}

func (js Jobs) String() string {
	var b strings.Builder
	for i, job := range js {
		b.WriteString(toLetter(job.id))
		if i != len(js)-1 {
			b.WriteString(", ")
		}
	}
	return b.String()
}

func toLetter(id int) string {
	if id >= 0 && id <= 26 {
		return string(rune('A' - 1 + id))
	}
	if id >= 27 && id <= 53 {
		return string(rune('a' - 27 + id))
	}
	return strconv.Itoa(id)
}
