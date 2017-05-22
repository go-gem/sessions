// Copyright 2016 The Gem Authors. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package sessions

import (
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

var (
	mutex sync.RWMutex
	data  = make(map[*fasthttp.RequestCtx]*Registry)
	datat = make(map[*fasthttp.RequestCtx]int64)
)

// Set stores a value in a given request.
func Set(ctx *fasthttp.RequestCtx, val *Registry) {
	mutex.Lock()
	data[ctx] = val
	datat[ctx] = time.Now().Unix()
	mutex.Unlock()
}

// Get returns a value stored in a given request.
func Get(ctx *fasthttp.RequestCtx) (val *Registry) {
	mutex.RLock()
	val = data[ctx]
	mutex.RUnlock()
	return
}

// GetOk returns stored value and presence state like multi-value return of map access.
func GetOk(ctx *fasthttp.RequestCtx) (*Registry, bool) {
	mutex.RLock()
	if v, ok := data[ctx]; ok {
		mutex.RUnlock()
		return v, ok
	}
	mutex.RUnlock()
	return nil, false
}

// Clear removes all values stored for a given request.
//
// It is no allow to keep the fasthttp.RequestCtx instance, so
// your application have to invoke Clear function at the end of
// a request lifetime.
//
// This is usually called by a handler wrapper to clean up request
// variables at the end of a request lifetime. See ClearHandler().
func Clear(ctx *fasthttp.RequestCtx) {
	mutex.Lock()
	clear(ctx)
	mutex.Unlock()
}

// clear is Clear without the lock.
func clear(ctx *fasthttp.RequestCtx) {
	if r, ok := data[ctx]; ok {
		r.close()
	}
	delete(data, ctx)
	delete(datat, ctx)
}

// ClearHandler wraps a fasthttp.RequestHandler and clears request values at the end
// of a request lifetime.
func ClearHandler(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		defer Clear(ctx)
		h(ctx)
	}
}
