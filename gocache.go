package gocache

import (
	"fmt"
	"time"
)

type DataObject struct {
	Data      string
	Timestamp time.Time
	Expiry    int64
}

var INMEM map[string]DataObject

func init() {
	INMEM = make(map[string]DataObject)
}
func SetData(key, data string, exp ...int64) {
	temp := DataObject{}
	temp.Data = data
	temp.Timestamp = time.Now()
	if len(exp) > 0 {
		temp.Expiry = exp[0]
	}
	INMEM[key] = temp
}

func GetData(key string) (string, bool) {
	retstr := ""
	retbool := false
	if v, ok := INMEM[key]; ok {
		retstr = v.Data
		retbool = true
		if INMEM[key].Expiry > 0 && time.Now().Sub(INMEM[key].Timestamp).Nanoseconds() >= INMEM[key].Expiry*1000000 {
			retstr = ""
			retbool = false
		}

	}
	return retstr, retbool
}

