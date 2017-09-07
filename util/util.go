package util

import "time"

func Assert(cond bool, msg string) {
	if !cond {
		panic(msg)
	}
}

func Tick() uint64 {
	return uint64(time.Now().UnixNano() / int64(time.Millisecond))
}
