package stride

import "testing"

func TestStrideNewJob(t *testing.T) {
	const (
		pass      = 0
		numerator = 10_000
		tickets   = 50
		stride    = numerator / tickets
	)
	j := NewJob(0, 50, t020)
	if j == nil {
		t.Fatal("stride.NewJob not implemented")
	}
	if j.Pass != pass || j.Tickets != tickets || j.Stride != stride {
		t.Errorf("NewJob(0, 50, 20ms) = (Pass=%d, Tickets=%d, Stride=%d), expected (Pass=%d, Tickets=%d, Stride=%d)", j.Pass, j.Tickets, j.Stride, pass, tickets, stride)
	}
}
