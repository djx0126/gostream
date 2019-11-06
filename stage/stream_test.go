package stage

import (
	"fmt"
	"strings"
	"testing"
)

func TestSliceStream(t *testing.T) {
	wordCount := make(map[string]int)

	words := []string{"a", "b", "A", "B", "c", "a"}
	StreamOf(words).Map(func(i interface{}) interface{} {
		str := i.(string)
		return strings.ToLower(str)
	}).Reduce(func(r interface{}, i interface{}) interface{} {
		str := i.(string)
		if _, ok := (wordCount)[str]; ok {
			(wordCount)[str]++
		} else {
			(wordCount)[str] = 1
		}
		return wordCount
	})
	fmt.Printf("end with: %v\n", wordCount)
}

func TestSlice(t *testing.T) {
	wordCount := make(map[string]int)

	newSourceStage(nil).Map(func(i interface{}) interface{} {
		str := i.(string)
		return strings.ToLower(str)
	}).Reduce(func(r interface{}, i interface{}) interface{} {
		str := i.(string)
		if _, ok := (wordCount)[str]; ok {
			(wordCount)[str]++
		} else {
			(wordCount)[str] = 1
		}
		return wordCount
	})
	fmt.Printf("end with: %v\n", wordCount)
}

