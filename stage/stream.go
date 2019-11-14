package stage

import (
	"errors"
	"reflect"
)

type ReducerFn func(interface{}, interface{}) interface{}

type MapFn func(interface{}) interface{}

type FilterFn func(interface{}) bool

func StreamOf(items interface{}) stage {
	switch reflect.TypeOf(items).Kind() {
	case reflect.Slice:
		return newSourceStage(newSliceContainer(items))
	default:
		errors.New("Unsupported type: " + reflect.TypeOf(items).Kind().String())
	}

	return newSourceStage(nil)
}
