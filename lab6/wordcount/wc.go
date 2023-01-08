package wordcount

import (
	"log"
	"os"
	"runtime"
	"unicode"
)

func loadMoby() []byte {
	moby, err := os.ReadFile("mobydick.txt")
	if err != nil {
		log.Fatal(err)
	}
	return moby
}

func wordCount(b []byte) (words int) {
	inWord := false
	for _, v := range b {
		r := rune(v)
		if unicode.IsSpace(r) && inWord {
			words++
		}
		inWord = unicode.IsLetter(r)
	}
	return
}

func shardSlice(input []byte, numShards int) (shards [][]byte) {
	shards = make([][]byte, numShards)
	if numShards < 2 {
		shards[0] = input[:]
		return
	}
	shardSize := len(input) / numShards
	start, end := 0, shardSize
	for i := 0; i < numShards; i++ {
		for j := end; j < len(input); j++ {
			char := rune(input[j])
			if unicode.IsSpace(char) {
				// split slice at position j, where there is a space
				// note: need to include the space in the shard to get accurate count
				end = j + 1
				shards[i] = input[start:end]
				start = end
				end += shardSize
				break
			}
		}
	}
	shards[numShards-1] = input[start:]
	return
}

func parallelWordCount(input []byte) (words int) {
	return doParallelWordCount(input, runtime.NumCPU())
}

func doParallelWordCount(input []byte, numShards int) (words int) {
	// TODO(student) implement parallel word count
	file := loadMoby()           //load the file
	wordc := wordCount(file)     //count hte words in the file
	shardSlice(input, numShards) //split provided byte into sub slices

	return wordc
}
