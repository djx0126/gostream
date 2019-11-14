package stage

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

func BenchmarkSimpleForLoop(b *testing.B) {
	s := prepareSlice()
	var x int

	b.ResetTimer()
	b.StartTimer()
	var itf interface{}
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(s); j++ {
			itf = s[j]
			x = itf.(int)
		}
	}
	b.StopTimer()
	x++
}

func BenchmarkAddressLoop(b *testing.B) {
	s := prepareSlice()
	var x int

	b.ResetTimer()
	b.StartTimer()
	var itf, itf1 interface{}
	itf = s
	for i := 0; i < b.N; i++ {
		efacePtr := (*eface)(unsafe.Pointer(&itf))
		typeItf := reflect.TypeOf(itf)
		typeEle := typeItf.Elem()
		eleSize := typeEle.Size()
		sliceDataPtr := (*slice)((*efacePtr).data)
		sliceSize := (*sliceDataPtr).len
		basePtr := uintptr(sliceDataPtr.array)

		for j := 0; j < sliceSize; j++ {
			u2intptr := basePtr + eleSize*uintptr(j)
			itf1 = toInterface(typeEle, u2intptr)
			x = itf1.(int)
		}
	}
	b.StopTimer()
	x++
}

func TestReflectReadSlice(t *testing.T) {
	s := prepareSlice()
	var itf, itf1 interface{}
	itf = s
	v := reflect.ValueOf(itf)

	itf1 = v.Index(4).Interface()
	x := itf1.(int)
	fmt.Printf("x=%d\n", x)
}

func BenchmarkReflectLoop(b *testing.B) {
	s := prepareSlice()
	var x int

	b.ResetTimer()
	b.StartTimer()
	var itf, itf1 interface{}
	itf = s
	for i := 0; i < b.N; i++ {
		v := reflect.ValueOf(itf)
		for j := 0; j < v.Len(); j++ {
			itf1 = v.Index(j).Interface()
			x = itf1.(int)
		}
	}
	b.StopTimer()
	x++
}

func prepareSlice() []int {
	s := make([]int, 100000)
	for j := 0; j < len(s); j++ {
		s[j] = j
	}
	return s
}

func TestBaseStage_CollectSlice(t *testing.T) {
	var ints []int

	fmt.Printf("eface.data %p, %v\n", &ints, uintptr(unsafe.Pointer(&ints)))
	collectToSlice(&ints, 11)
	ints = append(ints, 10)

	fmt.Printf("int %+v\n", ints)
}

func collectToSlice(slicePtr interface{}, x interface{}) {
	efacePtr := (*eface)(unsafe.Pointer(&slicePtr))

	fmt.Printf("eface.data %v\n", uintptr(efacePtr.data))
	slice := (*slice)(efacePtr.data)

	v := reflect.ValueOf(slicePtr).Elem()
	fmt.Println(v.Kind().String() + v.String())

	vx := reflect.ValueOf(x)
	v = reflect.Append(v, vx)

	(*slice).array = unsafe.Pointer(v.Pointer())
	(*slice).len = v.Len()
	(*slice).cap = v.Cap()
}
