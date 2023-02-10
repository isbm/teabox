package teaboxlib

import (
	"fmt"
	"strings"
	"sync"
)

// TeaboxRuntimeSession is a key/value storage per a module
type TeaboxRuntimeSession struct {
	kws map[string]interface{}
	mtx sync.RWMutex
}

// NewTeaboxRuntimeSession constructor
func NewTeaboxRuntimeSession() *TeaboxRuntimeSession {
	rts := new(TeaboxRuntimeSession)
	rts.kws = map[string]interface{}{}
	rts.mtx = *new(sync.RWMutex)

	return rts
}

// Set value to the session storage
func (rts *TeaboxRuntimeSession) Set(modname, k string, v interface{}) {
	rts.mtx.Lock()
	defer rts.mtx.Unlock()

	rts.kws[fmt.Sprintf("%s-%s", modname, k)] = v
}

// Get value from the session storage
func (rts *TeaboxRuntimeSession) Get(modname, k string) interface{} {
	rts.mtx.RLock()
	defer rts.mtx.RUnlock()

	v, ok := rts.kws[fmt.Sprintf("%s-%s", modname, k)]
	if !ok {
		return nil
	}

	return v
}

// Keys returns of all available keys of the session.
// It includes also public keys, prefixed with ":" colon.
func (rts *TeaboxRuntimeSession) Keys(modname string) []string {
	rts.mtx.RLock()
	defer rts.mtx.RUnlock()

	keys := []string{}
	for k := range rts.kws {
		if strings.HasPrefix(k, modname+"-") || strings.HasPrefix(k, ":") {
			keys = append(keys, k[len(modname):])
		}
	}

	return keys
}
func (rts *TeaboxRuntimeSession) Delete(modname, k string) {
	rts.mtx.Lock()
	defer rts.mtx.Unlock()

	delete(rts.kws, fmt.Sprintf("%s-%s", modname, k))
}

// Delete the entire session for the module. This does not affect
// public keys and data from other modules.
func (rts *TeaboxRuntimeSession) Flush(modname string) {
	keys := rts.Keys(modname)

	rts.mtx.Lock()
	defer rts.mtx.Unlock()

	for _, k := range keys {
		if strings.HasPrefix(k, ":") {
			continue
		}
		delete(rts.kws, k)
	}
}
