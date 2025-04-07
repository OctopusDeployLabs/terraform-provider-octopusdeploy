package internal

import (
	"log"
	"sync"
)

var Mutex = &sync.Mutex{}

var KeyedMutex = &MutexMap{}

type MutexMap struct {
	locks sync.Map
}

func (km *MutexMap) Lock(key string) {
	log.Printf("[DEBUG] KeyedMutex.Lock('%s') entry", key)
	mutex, found := km.locks.LoadOrStore(key, &sync.Mutex{})
	log.Printf("[DEBUG] KeyedMutex.Lock('%s') locking, exists:%t", key, found)
	mutex.(*sync.Mutex).Lock()
	log.Printf("[DEBUG] KeyedMutex.Lock('%s') exit", key)
}

func (km *MutexMap) Unlock(key string) {
	log.Printf("[DEBUG] KeyedMutex.Unlock('%s') entry", key)
	if mutex, found := km.locks.Load(key); found {
		log.Printf("[DEBUG] KeyedMutex.Unlock('%s') unlocking", key)
		mutex.(*sync.Mutex).Unlock()
	}
	log.Printf("[DEBUG] KeyedMutex.Unlock('%s') exit", key)
}
