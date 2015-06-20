package gocache

import (
	"sync"
	"time"
)

const (
	CLEANUP_DURATION = 60 * 30
)

type DataObject struct {
	Data      string
	Timestamp time.Time
	Expiry    int64
	lock      sync.Mutex
}

var inmem map[string]DataObject

func init() {
	inmem = make(map[string]DataObject)
	startCleanupTimer()
}

//SetData sets the data objecct of cache
//by initiating and linking it to the cache map
func SetData(key, data string, exp ...int64) {
	var temp DataObject
	if v, ok := inmem[key]; ok {
		temp = v
	} else {
		temp = DataObject{}
	}

	//make changes to object
	temp.lock.Lock()
	temp.Data = data
	temp.Timestamp = time.Now()
	if len(exp) > 0 {
		temp.Expiry = exp[0]
	}
	temp.lock.Unlock()

	//update cache map
	inmem[key] = temp
}

//GetData gets from cache
func GetData(key string) (string, bool) {
	retstr := ""
	retbool := false
	if v, ok := inmem[key]; ok {
		retstr = v.Data
		retbool = true
		if checkExpiry(key) {
			retstr = ""
			retbool = false
		}

	}
	return retstr, retbool
}

func cleanup() {
	for key := range inmem {
		if checkExpiry(key) {
			delete(inmem, key)
		}
	}
}

func startCleanupTimer() {
	ticker := time.Tick(CLEANUP_DURATION)
	go (func() {
		for {
			select {
			case <-ticker:
				cleanup()
			}
		}
	})()
}

//returns true if a key has expired
func checkExpiry(key string) bool {
	// return false if key not found
	if _, ok := inmem[key]; ok {
		return false
	}
	val := inmem[key]
	//return false if ttl not set
	ttlSecs := val.Expiry * 1000000
	if ttlSecs <= 0 {
		return false
	}

	now := time.Now()
	keyEntryTime := val.Timestamp
	timeSinceEntry := now.Sub(keyEntryTime).Nanoseconds()

	// return true if its been longer than ttl
	return timeSinceEntry >= ttlSecs
}
