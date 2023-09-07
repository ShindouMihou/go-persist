# go-persist

a little persistence utility for golang applications when all you need is just some bit of persisted data, and not an 
entire database. go-persist is that exactly with opt-able conveniences such as auto-persist, array persistence and 
even map persistence.

## Demo

### Array Persistence
```go
func main() {
	wd, _ := os.Getwd()
	array := persist.NewArray[string](filepath.Join(wd, "data", "array.json"))
    	if err := array.Load(); err != nil {
		log.Panicln("failed to load array", err)
    	}
	array.Append("hello", "world")
	
	// We recommend turning on auto-persistence after loading the already-persisted data
	// because we may end up overwriting the data, although highly unlikely as the internal backing 
	// does not overwrite unless there has been a change in the value.
    	array.Persist(5 * time.Second)
	
	// You can also opt to manually flushing the new changes.
	if err := array.Save(); err != nil {
        log.Panicln("failed to save array", err)
    }
}
```


### Set Persistence
this is backed with [deckarep/golang-set](https://github.com/deckarep/golang-set).

```go
func main() {
	wd, _ := os.Getwd()
	set := persist.NewSet[string](filepath.Join(wd, "data", "set.json"))
    	if err := set.Load(); err != nil {
		log.Panicln("failed to load set", err)
    	}
	set.Append("hello", "world")
	
	if set.Contains("world") {
	    set.Append("and", "galaxy")	
    	}
	
	// We recommend turning on auto-persistence after loading the already-persisted data
	// because we may end up overwriting the data, although highly unlikely as the internal backing 
	// does not overwrite unless there has been a change in the value.
    	set.Persist(5 * time.Second)
	
	// You can also opt to manually flushing the new changes.
	if err := set.Save(); err != nil {
        	log.Panicln("failed to save set", err)
    	}
}
```

### Map Persistence
```go
func main() {
	wd, _ := os.Getwd()
	pmap := persist.NewMap[string, string](filepath.Join(wd, "data", "map.json"))
    	if err := map.Load(); err != nil {
		log.Panicln("failed to save map", err)
    	}
    	pmap.Set("hello", "world")
	
	// We recommend turning on auto-persistence after loading the already-persisted data
	// because we may end up overwriting the data, although highly unlikely as the internal backing 
	// does not overwrite unless there has been a change in the value.
	pmap.Persist(5 * time.Second)
	
	// You can also opt to manually flushing the new changes.
	if err := pmap.Save(); err != nil {
        	log.Panicln("failed to save map", err)
    	}
}
```

### Interface Persistence

```go
type Hello {
	World string `json:"hello"`
}

func main() {
	wd, _ := os.Getwd()
	ptype := persist.NewPersisted[Hello](filepath.Join(wd, "data", ".json"), nil)
    	if err := ptype.Load(); err != nil {
		log.Panicln("failed to save type", err)
    	}
	ptype.Set(&Hello{World:"world"})
	
	// You can also edit the value directly.
	ptype.Edit(func (value *Hello) {
	    value.World = "galaxy"	
    	})
	
	// We recommend turning on auto-persistence after loading the already-persisted data
	// because we may end up overwriting the data, although highly unlikely as the internal backing 
	// does not overwrite unless there has been a change in the value.
    	ptype.Persist(5 * time.Second)
	
	// You can also opt to manually flushing the new changes.
	if err := ptype.Save(); err != nil {
        	log.Panicln("failed to save type", err)
    	}
}
```
