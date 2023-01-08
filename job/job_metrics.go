package job

import (
	"dat320/lab4/scheduler/system/systime"
	"time"
)

func (j *Job) Scheduled(s systime.SystemTime) {
	//j.SystemTime = s
	j.SystemTime = s
	j.arrival = j.SystemTime.Now()
	// (student) implement task 2.1
}
func (j *Job) Started(cpuID int) {
	// (student) implement task 2.2
	//j.id = cpuID
	if j.start == NotStartedYet {
		j.start = j.SystemTime.Now() //for fifo
	}
}
func (j Job) TurnaroundTime() time.Duration {
	r := j.finished - j.arrival
	return r

}
func (j Job) ResponseTime() time.Duration {
	// (student) implement task 2.3
	r := j.start - j.arrival //fifo
	return r
}
