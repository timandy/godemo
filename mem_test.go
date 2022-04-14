package main

import (
	"github.com/stretchr/testify/assert"
	"runtime/debug"
	"sync"
	"testing"
	"unsafe"
)

type ExtendedMap struct {
	labels map[string]string
	flag   int
}

//go:norace
//go:nocheckptr
func TestForMap(t *testing.T) {
	const concurrency = 5000
	const loopTimes = 1000
	wg := &sync.WaitGroup{}
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			var ptr unsafe.Pointer
			for j := 0; j < loopTimes; j++ {
				m := make(map[string]string)
				m["key"] = "value"
				ptr = unsafe.Pointer(&m)

				assert.NotNil(t, ptr)
				_, flag := getEMap(ptr)
				assert.NotEqual(t, 12345678, flag)

				labels := make(map[string]string)
				labels["key"] = "value"
				em := ExtendedMap{labels: labels, flag: 12345678}
				ptr = unsafe.Pointer(&em)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

//go:norace
//go:nocheckptr
func getEMap(ptr unsafe.Pointer) (em *ExtendedMap, flag int) {
	panicOnFault := debug.SetPanicOnFault(true)
	defer func() {
		debug.SetPanicOnFault(panicOnFault)
		recover()
	}()
	em = (*ExtendedMap)(ptr)
	return em, em.flag
}
