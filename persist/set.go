package persist

import (
	ptr "github.com/ShindouMihou/go-little-utils"
	mapset "github.com/deckarep/golang-set/v2"
	"time"
)

type Set[Value comparable] struct {
	backing *Persisted[mapset.Set[Value]]
}

func NewSet[Value comparable](path string) *Set[Value] {
	return &Set[Value]{
		backing: NewPersisted[mapset.Set[Value]](path, ptr.Ptr(mapset.NewSet[Value]())),
	}
}

func (pset *Set[Value]) Pop() (Value, bool) {
	return pset.backing.Get().Pop()
}

func (pset *Set[Value]) Contains(values ...Value) bool {
	return pset.backing.Get().Contains(values...)
}

func (pset *Set[Value]) ContainsAny(values ...Value) bool {
	return pset.backing.Get().ContainsAny(values...)
}

func (pset *Set[Value]) Append(values ...Value) {
	pset.backing.Edit(func(val *mapset.Set[Value]) {
		(*val).Append(values...)
	})
}

func (pset *Set[Value]) Each(fn func(value Value) bool) {
	pset.backing.Get().Each(fn)
}

func (pset *Set[Value]) Persist(every time.Duration) {
	pset.backing.Persist(every)
}

func (pset *Set[Value]) Watch() (close func() error, err error) {
	return pset.backing.Watch()
}

func (pset *Set[Value]) Save() error {
	return pset.backing.Save()
}

func (pset *Set[Value]) Load() error {
	return pset.backing.Load()
}

func (pset *Set[Value]) Length() int {
	return pset.backing.Get().Cardinality()
}
