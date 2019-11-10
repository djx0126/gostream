package stage

import (
	"fmt"
	"strings"
	"testing"
)

func TestWordCount(t *testing.T) {
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
	wordCount := StreamOf(words).Limit(5).Map(func(i interface{}) interface{} {
		return strings.ToLower(i.(string))
	}).Reduce(func(r interface{}, i interface{}) interface{} {
		str := i.(string)
		if r == nil {
			r = make(map[string]int)
		}

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

func TestSliceStreamFilter(t *testing.T) {
	words := []string{"a", "b", "A", "B", "c", "a"}
	wordCount := StreamOf(words).Filter(func(i interface{}) bool {
		return strings.Compare(i.(string), "B") != 0
	}).Map(func(i interface{}) interface{} {
		return strings.ToLower(i.(string))
	}).Reduce(func(r interface{}, i interface{}) interface{} {
		str := i.(string)
		if r == nil {
			r = make(map[string]int)
		}

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

func TestNewSliceCollector(t *testing.T) {
	words := []string{"a", "b", "A", "B", "c", "a"}

	var limitWords []string
	StreamOf(words).Limit(3).Collect(NewSliceCollector(&limitWords))

	fmt.Printf("limit words: %v\n", limitWords)

	var noWords []string
	StreamOf(words).Filter(func(i interface{}) bool {
		return strings.Compare(i.(string), "no") == 0
	}).Collect(NewSliceCollector(&noWords))

	fmt.Printf("no words: %v\n", noWords)
}
