package stage

import (
	"errors"
	"reflect"
	"unsafe"
)

type container interface {
	newIterator() iterator
}

type iterator interface {
	next() interface{}
	hasNext() bool
}

type slice struct {
	array unsafe.Pointer
	len   int
	cap   int
}
type eface struct {
	_typePtr *_type
	data     unsafe.Pointer
}

type tflag uint8
type typeAlg struct {
	// function for hashing objects of this type
	// (ptr to object, seed) -> hash
	hash func(unsafe.Pointer, uintptr) uintptr
	// function for comparing objects of this type
	// (ptr to object A, ptr to object B) -> ==?
	equal func(unsafe.Pointer, unsafe.Pointer) bool
}
type nameOff int32
type typeOff int32
type _type struct {
	size       uintptr
	ptrdata    uintptr // size of memory prefix holding all pointers
	hash       uint32
	tflag      tflag
	align      uint8
	fieldalign uint8
	kind       uint8
	alg        *typeAlg
	// gcdata stores the GC type data for the garbage collector.
	// If the KindGCProg bit is set in kind, gcdata is a GC program.
	// Otherwise it is a ptrmask bitmap. See mbitmap.go for details.
	gcdata    *byte
	str       nameOff
	ptrToThis typeOff
}

func toInterface(typeEle reflect.Type, ptr uintptr) interface{} {
	typeDataPtr:=(*(*eface)(unsafe.Pointer(&typeEle))).data
	var i interface{}
	e := (*eface)(unsafe.Pointer(&i))
	e.data = unsafe.Pointer(ptr)
	tp := (*_type)(typeDataPtr)
	e._typePtr = tp
	return i
}

type sliceContainer struct {
	efacePtr *eface
	eleSize uintptr
	typeItf reflect.Type
	typeEle reflect.Type
	sliceSize int
	basePtr uintptr
}

func newSliceContainer(items interface{}) container {
	s := sliceContainer{}

	s.efacePtr = (*eface)(unsafe.Pointer(&items))
	s.typeItf = reflect.TypeOf(items)

	switch s.typeItf.Kind() {
	case reflect.Slice:
		s.typeEle = s.typeItf.Elem()
		s.eleSize = s.typeEle.Size()
	default:
		errors.New("Unsupported type: " + s.typeItf.String())
	}

	sliceDataPtr := (*slice)((*s.efacePtr).data)
	s.sliceSize = (*sliceDataPtr).len
	s.basePtr = uintptr(sliceDataPtr.array)

	return &s
}

type sliceContainerIterator struct {
	container *sliceContainer
	index int
}

func (s *sliceContainer) newIterator()  iterator {
	return &sliceContainerIterator{
		container:s,
	}
}

func (iter *sliceContainerIterator) hasNext() bool {
	return iter.index < iter.container.sliceSize
}

func (iter *sliceContainerIterator) next()  (r interface{}) {
	u2intptr := iter.container.basePtr + iter.container.typeEle.Size()*uintptr(iter.index)
	r = toInterface(iter.container.typeEle, u2intptr)
	iter.index++
	return
}



