package persist

import (
	ptr "github.com/ShindouMihou/go-little-utils"
	"time"
)

type Map[Key comparable, Value any] struct {
	backing *Persisted[map[Key]Value]
}

func NewMap[Key comparable, Value any](path string) *Map[Key, Value] {
	return &Map[Key, Value]{
		backing: NewPersisted[map[Key]Value](path, ptr.Ptr(make(map[Key]Value))),
	}
}

func (pmap *Map[Key, Value]) Get(key Key) *Value {
	return pmap.GetOr(key, nil)
}

func (pmap *Map[Key, Value]) Each(fn func(key Key, value Value)) {
	backing := pmap.backing.Get()
	for key, value := range backing {
		fn(key, value)
	}
}

func (pmap *Map[Key, Value]) Length() int {
	return len(pmap.backing.Get())
}

func (pmap *Map[Key, Value]) GetOr(key Key, def *Value) *Value {
	val, ok := pmap.backing.Get()[key]
	if !ok {
		return def
	}
	return &val
}

func (pmap *Map[Key, Value]) Set(key Key, value Value) {
	pmap.backing.Edit(func(val *map[Key]Value) {
		(*val)[key] = value
	})
}

func (pmap *Map[Key, Value]) Persist(every time.Duration) {
	pmap.backing.Persist(every)
}

func (pmap *Map[Key, Value]) Save() error {
	return pmap.backing.Save()
}

func (pmap *Map[Key, Value]) Load() error {
	return pmap.backing.Load()
}
