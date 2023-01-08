package wordcount

import (
	"testing"
)

const expectedWC = 175131

var mobyBytes []byte

func init() {
	mobyBytes = loadMoby()
}

func TestWordCount(t *testing.T) {
	cnt := wordCount(mobyBytes)
	if cnt != expectedWC {
		t.Errorf("expected %d words, found %d words\n", expectedWC, cnt)
	}
}

func TestParallelWordCount(t *testing.T) {
	cnt := parallelWordCount(mobyBytes)
	if cnt != expectedWC {
		t.Errorf("expected %d words, found %d words\n", expectedWC, cnt)
	}
}

func TestParallelWordCountManyShards(t *testing.T) {
	wc := wordCount(mobyBytes)
	for shards := 1; shards < 10; shards++ {
		pwc := doParallelWordCount(mobyBytes, shards)
		if pwc != wc {
			t.Errorf("doParallelWordCount(..., %d)=%d, expected %d", shards, pwc, wc)
		}
	}
}

var result int

func BenchmarkWordCountSequential(b *testing.B) {
	var wc int
	for i := 0; i < b.N; i++ {
		wc = wordCount(mobyBytes)
	}
	result = wc
}

func BenchmarkWordCountParallel(b *testing.B) {
	var wc int
	for i := 0; i < b.N; i++ {
		wc = parallelWordCount(mobyBytes)
	}
	result = wc
}
