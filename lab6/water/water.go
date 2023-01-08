package water

import (
	"sync"
)

// Water structure holds the synchronization primitives and
// data required to solve the water molecule problem.
// moleculeCount holds the number of molecules formed so far.
// result string contains the sequence of "H" and "O".
// wg WaitGroup is used to wait for goroutine completion.
type Water struct {
	wg            sync.WaitGroup
	moleculeCount int
	result        string
	// TODO(student) add missing fields, if necessary
}

// New initializes the water structure.
func New() *Water {
	water := &Water{}
	// TODO(student) initialize the Water struct
	return water
}

// releaseOxygen produces one oxygen atom if no oxygen atom is already present.
// If an oxygen atom is already present, it will block until enough hydrogen
// atoms have been produced to consume the atoms necessary to produce water.
//
// The w.wg.Done() must be called to indicate the completion of the goroutine.
func (w *Water) releaseOxygen() {
	defer w.wg.Done()
	// TODO(student) implement the releaseOxygen routine
}

// releaseHydrogen produces one hydrogen atom unless two hydrogen atoms are already present.
// If two hydrogen atoms are already present, it will block until another oxygen
// atom has been produced to consume the atoms necessary to produce water.
//
// The w.wg.Done() must be called to indicate the completion of the goroutine.
func (w *Water) releaseHydrogen() {
	defer w.wg.Done()
	// TODO(student) implement the releaseHydrogen routine
}

// produceMolecule forms the water molecules.
func (w *Water) produceMolecule(done chan bool) {
	// TODO(student) implement the produceMolecule routine
	done <- true
}

func (w *Water) finish() {
	// TODO(student) implement the finish routine to complete the water molecule formation
}

// Molecules returns the number of water molecules that has been created.
func (w *Water) Molecules() int {
	// TODO(student) Add any missing code
	return w.moleculeCount
}

// Make returns a sequence of water molecules derived from the input of hydrogen and oxygen atoms.
// DO NOT edit the Make method or modify the signatures of the other methods used by Make.
func (w *Water) Make(input string) string {
	done := make(chan bool)
	go w.produceMolecule(done)
	for _, ch := range input {
		w.wg.Add(1)
		switch ch {
		case 'O':
			go w.releaseOxygen()
		case 'H':
			go w.releaseHydrogen()
		}
	}
	w.wg.Wait()
	w.finish()
	<-done
	return w.result
}
