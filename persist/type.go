package persist

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ShindouMihou/chikador/chikador"
	"github.com/ShindouMihou/go-little-utils/fileutils"
	"github.com/ShindouMihou/go-little-utils/locks"
	"github.com/fsnotify/fsnotify"
	"io"
	"os"
	"sync"
	"time"
)

type Persisted[T any] struct {
	file            string
	value           *T
	hasBeenModified bool
	mutex           sync.RWMutex
	chismis         *chikador.Chismis
	filechanged     bool
}

func NewPersisted[T any](path string, initial *T) *Persisted[T] {
	return &Persisted[T]{
		file:            path,
		value:           initial,
		hasBeenModified: false,
		mutex:           sync.RWMutex{},
	}
}

func (persisted *Persisted[T]) Load() error {
	file, err := os.Open(persisted.file)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("failed to close backing file: ", err)
		}
	}(file)

	bytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	var t T
	if err := json.Unmarshal(bytes, &t); err != nil {
		return err
	}
	persisted.value = &t
	return nil
}

func (persisted *Persisted[T]) Set(value T) {
	persisted.hasBeenModified = true
	locks.UseWrite(&persisted.mutex, func() {
		persisted.value = &value
	})
}

func (persisted *Persisted[T]) Get() T {
	persisted.mutex.RLock()
	defer persisted.mutex.RUnlock()
	return *persisted.value
}

func (persisted *Persisted[T]) Edit(fn func(val *T)) {
	persisted.hasBeenModified = true
	locks.UseWrite(&persisted.mutex, func() {
		fn(persisted.value)
	})
}

func (persisted *Persisted[T]) Save() error {
	if persisted.value == nil {
		return nil
	}
	res, err := json.Marshal(persisted.value)
	if err != nil {
		return err
	}
	return fileutils.SaveOrOverwrite(persisted.file, res)
}

func (persisted *Persisted[T]) Persist(every time.Duration) {
	go persisted.persist(every)
}

func (persisted *Persisted[T]) Watch() (close func() error, err error) {
	if persisted.chismis != nil {
		return nil, errors.New("already watching for changes in the file")
	}
	if err := persisted.Save(); err != nil {
		return nil, err
	}
	chismis, err := chikador.Watch(persisted.file, chikador.WithDedupe)
	if err != nil {
		return nil, err
	}
	persisted.chismis = chismis
	persisted.chismis.Listen(func(msg *chikador.Message) {
		if !msg.Event.Has(fsnotify.Write) {
			return
		}
		// only change when it's not a change coming from us
		if persisted.filechanged {
			persisted.filechanged = false
			return
		}
		if err := persisted.Load(); err != nil {
			fmt.Println("err(go-persist): failed to reload from ", persisted.file, ":", err)
		}
		fmt.Println("reloaded with value: ", *persisted.value)
	})
	return func() error {
		err := persisted.chismis.Close()
		persisted.chismis = nil
		return err
	}, nil
}

func (persisted *Persisted[T]) persist(every time.Duration) {
	time.Sleep(every)

	persisted.mutex.RLock()
	defer persisted.mutex.RUnlock()
	defer persisted.persist(every)

	if !persisted.hasBeenModified {
		return
	}

	persisted.filechanged = true

	if err := persisted.Save(); err != nil {
		fmt.Println("failed to save persistent data: ", err)
		return
	}

	persisted.hasBeenModified = false
}
