package paging

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPTAppend(t *testing.T) {
	for i, test := range PTAppendTests {
		pageTable := PageTable{test.pageTable}
		pageTable.Append(test.in)

		if diff := cmp.Diff(test.want, pageTable.frameIndices); diff != "" {
			t.Errorf("PTAppendTests %d: Unexpected page table content; (-want +got):\n%s", i, diff)
		}
	}
}

func TestPTFree(t *testing.T) {
	for i, test := range PTFreeTests {
		pageTable := PageTable{test.pageTable}
		freed, err := pageTable.Free(test.in)

		if test.want.err == nil && err != nil {
			t.Errorf("PTFreeTests %d: Failed to free pages when it was expected to succeed; want no error, got %v", i, err)
		}

		if test.want.err != nil && err == nil {
			t.Errorf("PTFreeTests %d: Successfully freed pages when it was expected to fail", i)
		}

		if diff := cmp.Diff(test.want.freed, freed); diff != "" {
			t.Errorf("PTFreeTests %d: Unexpected list of freed pages; (-want +got):\n%s", i, diff)
		}

		if diff := cmp.Diff(test.want.pageTable, pageTable.frameIndices); diff != "" {
			t.Errorf("PTFreeTests %d: Unexpected page table content after free operation; (-want +got):\n%s", i, diff)
		}
	}
}

func TestPTLookup(t *testing.T) {
	for i, test := range PTLookupTests {
		pageTable := PageTable{test.pageTable}
		frameIndex, err := pageTable.Lookup(test.in)

		if test.want.err == nil && err != nil {
			t.Errorf("PTFreeTests %d: Failed to lookup address %d when it was expected to succeed; want no error, got %v", i, test.in, err)
		}

		if test.want.err != nil && err == nil {
			t.Errorf("PTFreeTests %d: Successfully looked up address %d when it was expected to fail", i, test.in)
		}

		if diff := cmp.Diff(test.want.frameIndex, frameIndex); diff != "" {
			t.Errorf("PTFreeTests %d: Got invalid page number from lookup; (-want +got):\n%s", i, diff)
		}
	}
}
