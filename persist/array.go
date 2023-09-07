package persist

import (
	ptr "github.com/ShindouMihou/go-little-utils"
	"time"
)

type Array[Value any] struct {
	backing *Persisted[[]Value]
}

func NewArray[Value any](path string) *Array[Value] {
	return &Array[Value]{
		backing: NewPersisted[[]Value](path, ptr.Ptr([]Value{})),
	}
}

func (parr *Array[Value]) Get(index int) *Value {
	return parr.GetOr(index, nil)
}

func (parr *Array[Value]) GetOr(index int, def *Value) *Value {
	arr := parr.backing.Get()
	if len(arr) < (index + 1) {
		return def
	}
	val := arr[index]
	return &val
}

func (parr *Array[Value]) Each(fn func(index int, value Value)) {
	backing := parr.backing.Get()
	for index, value := range backing {
		fn(index, value)
	}
}

func (parr *Array[Value]) Contains(predicate func(value Value) bool) bool {
	backing := parr.backing.Get()
	for _, value := range backing {
		if predicate(value) {
			return true
		}
	}
	return false
}

// Set sets the value at that specified index, if the array doesn't reach that far yet
// then it will ignore.
func (parr *Array[Value]) Set(index int, value Value) {
	parr.backing.Edit(func(val *[]Value) {
		if len(*val) < (index + 1) {
			return
		}
		(*val)[index] = value
	})
}

// UnsafeSet sets the value at that specified index, if the array doesn't reach that far, then
// it will throw a panic.
func (parr *Array[Value]) UnsafeSet(index int, value Value) {
	parr.backing.Edit(func(val *[]Value) {
		(*val)[index] = value
	})
}

func (parr *Array[Value]) Append(values ...Value) {
	parr.backing.Edit(func(val *[]Value) {
		*val = append(*val, values...)
	})
}

func (parr *Array[Value]) Persist(every time.Duration) {
	parr.backing.Persist(every)
}

func (parr *Array[Value]) Save() error {
	return parr.backing.Save()
}

func (parr *Array[Value]) Load() error {
	return parr.backing.Load()
}

func (parr *Array[Value]) Length() int {
	return len(parr.backing.Get())
}
