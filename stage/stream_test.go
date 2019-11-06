package stage

import (
	"fmt"
	"strings"
	"testing"
)

func TestWordCount(t *testing.T)  {
	words := []string{"a", "b", "A", "B", "c", "a"}
	counter := make(map[string]int)

	for _, word := range words {
		str := strings.ToLower(word)
		if _, ok := (counter)[str]; ok {
			(counter)[str]++
		} else {
			(counter)[str] = 1
		}
	}
	fmt.Printf("end with: %v\n", counter)
}

func TestSliceStream(t *testing.T) {
	words := []string{"a", "b", "A", "B", "c", "a"}
	wordCount := StreamOf(words).Map(func(i interface{}) interface{} {
		return strings.ToLower(i.(string))
	}).Reduce(func(r interface{}, i interface{}) interface{} {
		str := i.(string)
		if r == nil { r = make(map[string]int)}

		counter := r.(map[string]int)
		if _, ok := (counter)[str]; ok {
			(counter)[str]++
		} else {
			(counter)[str] = 1
		}
		return counter
	}).(map[string]int)
	fmt.Printf("end with: %v\n", wordCount)
}

func TestSlice(t *testing.T) {
	wordCount := make(map[string]int)

	newSourceStage(nil).Map(func(i interface{}) interface{} {
		return strings.ToLower(i.(string))
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

