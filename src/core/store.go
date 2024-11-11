package core

import (
	"sync"
	"time"
)

// set mutex
var StoreObj Store = Store{entries: map[Key]Value{}}

type Store struct {
	timer   int64
	entries map[Key]Value
}

type Key struct {
	key     interface{}
	keyType string
}

type Value struct {
	value   interface{}
	lock    *sync.Mutex
	valType string
	expiry  int64
}

func (this *Store) Set(key interface{}, val interface{}, ttl int64) (bool, error) {
	var expiry int64 = -1
	if ttl != -1 {
		expiry = time.Now().UnixMilli() + ttl
	}
	keyWr := Key{key: key}
	valWr, isKeyPresent := this.entries[keyWr]
	if isKeyPresent {
		valWr.lock.Lock()
		valWr.value = val
		valWr.expiry = expiry
		valWr.lock.Unlock()
	} else {
		valWr = Value{value: val, expiry: expiry, lock: &sync.Mutex{}}
	}

	this.entries[keyWr] = valWr
	return true, nil
}

func (this *Store) Get(key interface{}) (interface{}, error) {
	keyWr := Key{key: key}
	valWr, isPresent := this.entries[keyWr]
	if !isPresent {
		return "-1", nil
	}

	valExpiry, timeNow := valWr.expiry, time.Now().UnixMilli()
	if valExpiry != -1 && valExpiry < timeNow {
		delete(this.entries, keyWr)
		return "-1", nil
	}
	return valWr.value, nil
}
