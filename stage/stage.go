package stage

import (
	"fmt"
)

type stage interface {
	Reduce(rf ReducerFn) interface{}
	Map(mf MapFn) stage
}

type baseStage struct {
	sourceStage   *sourceStage
	previousStage *baseStage
	nextStage     *baseStage
	sink          chan interface{}
	acceptFn      func(interface{})
}

type sourceStage struct {
	baseStage
	container container
}

type actionStage struct {
	baseStage
}

func (s *baseStage) Reduce(rf ReducerFn) (rt interface{}) {
	stage := newActionStage(s.sourceStage, s)
	var r interface{}
	stage.acceptFn = func(i interface{}) {
		r = rf(r, i)
	}
	stage.evaluate()
	return r
}

func (s *baseStage) Map(mf MapFn) stage {
	stage := newStage(s.sourceStage, s)
	stage.acceptFn = func(i interface{}) {
		//stage.sink <- mf(i)
		stage.nextStage.acceptFn(mf(i))
	}
	return stage
}

func (s *sourceStage) iterate() {
	var items []interface{}
	if s.container == nil {
		fmt.Printf("use mock data!\n")
		items = []interface{}{"1a", "2b", "3A", "4B", "5c", "6a"}
	} else {
		items = s.container.contents()
	}
	for _, item := range items {
		//s.sink <- item
		s.nextStage.acceptFn(item)
	}
}

func (s *actionStage) evaluate() interface{} {
	s.sourceStage.iterate()
	return nil
}

func newSourceStage(container container) stage {
	stage := &sourceStage{baseStage{
		sourceStage:   nil,
		previousStage: nil,
		sink:          make(chan interface{}),
	}, container,
	}
	stage.sourceStage = stage

	return stage
}

func newActionStage(sourceStage *sourceStage, previousStage *baseStage) *actionStage {
	stage := &actionStage{
		baseStage{sourceStage: sourceStage, previousStage: previousStage},
	}
	previousStage.nextStage = &stage.baseStage

	stage.attachChan()
	return stage
}

func newStage(sourceStage *sourceStage, previousStage *baseStage) *baseStage {
	stage := &baseStage{
		sourceStage:   sourceStage,
		previousStage: previousStage,
		sink:          make(chan interface{}),
	}
	previousStage.nextStage = stage

	// attach chan to previous stage
	stage.attachChan()
	return stage
}

func (s *baseStage) attachChan() {
	go func() {
		for item := range s.previousStage.sink {
			s.acceptFn(item)
		}
	}()
}

