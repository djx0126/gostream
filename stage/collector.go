package stage

import (
	"reflect"
	"unsafe"
)

type collector interface {
	accept(x interface{})
}

type sliceCollector struct {
	slicePtr interface{}
	sliceValue reflect.Value
	_sliceP *slice
}

func (c *sliceCollector) accept(x interface{}) {
	c.sliceValue = reflect.Append(c.sliceValue, reflect.ValueOf(x))

	(*c._sliceP).len = c.sliceValue.Len()
	(*c._sliceP).cap = c.sliceValue.Cap()
}

func NewSliceCollector(slicePtr interface{}) collector {
	efacePtr := (*eface)(unsafe.Pointer(&slicePtr))

	slice := (*slice)(efacePtr.data)

	v := reflect.ValueOf(slicePtr).Elem()
	(*slice).array = unsafe.Pointer(v.Pointer())
	return &sliceCollector{
		slicePtr: slicePtr,
		sliceValue:v,
		_sliceP:slice,
	}
}
