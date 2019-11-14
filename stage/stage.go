package stage

import "reflect"

type stage interface {
	Reduce(rf ReducerFn) interface{}
	Map(mf MapFn) stage
	MapS(mf interface{}) stage
	Filter(fn FilterFn) stage
	Limit(l int) stage
	Collect(c collector)
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

func (s *baseStage) Reduce(rf ReducerFn) interface{} {
	stage := newActionStage(s.sourceStage, s)
	var r interface{}
	stage.acceptFn = func(i interface{}) {
		r = rf(r, i)
	}
	stage.evaluate()
	return r
}

func (s *baseStage) Collect(c collector) {
	stage := newActionStage(s.sourceStage, s)
	stage.acceptFn = func(i interface{}) {
		c.accept(i)
	}
	stage.evaluate()
}

func (s *baseStage) Map(mf MapFn) stage {
	stage := newStage(s.sourceStage, s)
	stage.acceptFn = func(i interface{}) {
		//stage.sink <- mf(i)
		stage.nextStage.acceptFn(mf(i))
	}
	return stage
}

func (s *baseStage) MapS(fn interface{}) stage {
	fv := reflect.ValueOf(fn)
	input := make([]reflect.Value, 1)

	stage := newStage(s.sourceStage, s)
	stage.acceptFn = func(i interface{}) {
		//stage.sink <- mf(i)
		input[0] = reflect.ValueOf(i)
		stage.nextStage.acceptFn(fv.Call(input)[0].Interface())
	}
	return stage
}

func (s *baseStage) Filter(filter FilterFn) stage {
	stage := newStage(s.sourceStage, s)
	stage.acceptFn = func(i interface{}) {
		if filter(i) {
			//stage.sink <- mf(i)
			stage.nextStage.acceptFn(i)
		}
	}
	return stage
}

func (s *baseStage) Limit(l int) stage {
	stage := newStage(s.sourceStage, s)
	accepted := 0
	stage.acceptFn = func(i interface{}) {
		if accepted < l {
			accepted++
			//stage.sink <- mf(i)
			stage.nextStage.acceptFn(i)
		}

	}
	return stage
}

func (s *sourceStage) iterate() {
	for iterator := s.container.newIterator(); iterator.hasNext(); {
		item := iterator.next()
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
