// Copyright 2016 The Gem Authors. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package sessions

import (
	"testing"

	"github.com/valyala/fasthttp"
)

func TestContext(t *testing.T) {
	assertEqual := func(val, exp *Registry) {
		if val != exp {
			t.Errorf("Expected %v, got %v.", exp, val)
			t.FailNow()
		}
	}

	ctx := &fasthttp.RequestCtx{}
	ctx2 := &fasthttp.RequestCtx{}

	registry := &Registry{}

	// Get()
	assertEqual(Get(ctx), nil)

	// Set()
	Set(ctx, registry)
	assertEqual(Get(ctx), registry)
	if len(data) != 1 {
		t.Errorf("Expected %v, got %v.", 1, len(data))
		t.FailNow()
	}
	if len(datat) != 1 {
		t.Errorf("Expected %v, got %v.", 1, len(datat))
		t.FailNow()
	}

	//GetOk
	value, ok := GetOk(ctx)
	assertEqual(value, registry)
	if ok != true {
		t.Errorf("Expected %v, got %v.", true, ok)
		t.FailNow()
	}

	value, ok = GetOk(ctx2)
	assertEqual(value, nil)
	if ok != false {
		t.Errorf("Expected %v, got %v.", false, ok)
		t.FailNow()
	}

	Set(ctx2, nil)
	value, ok = GetOk(ctx2)
	assertEqual(value, nil)
	if ok != true {
		t.Errorf("Expected %v, got %v.", true, ok)
		t.FailNow()
	}

	// Clear()
	Clear(ctx)
	value, ok = GetOk(ctx)
	assertEqual(value, nil)
	if ok != false {
		t.Errorf("Expected %v, got %v.", false, ok)
		t.FailNow()
	}
}

func parallelReader(ctx *fasthttp.RequestCtx, iterations int, wait, done chan struct{}) {
	<-wait
	for i := 0; i < iterations; i++ {
		Get(ctx)
	}
	done <- struct{}{}

}

func parallelWriter(ctx *fasthttp.RequestCtx, value *Registry, iterations int, wait, done chan struct{}) {
	<-wait
	for i := 0; i < iterations; i++ {
		Set(ctx, value)
	}
	done <- struct{}{}

}

func benchmarkMutex(b *testing.B, numReaders, numWriters, iterations int) {
	b.StopTimer()
	ctx := &fasthttp.RequestCtx{}
	registry := &Registry{}
	done := make(chan struct{})
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		wait := make(chan struct{})

		for i := 0; i < numReaders; i++ {
			go parallelReader(ctx, iterations, wait, done)
		}

		for i := 0; i < numWriters; i++ {
			go parallelWriter(ctx, registry, iterations, wait, done)
		}

		close(wait)

		for i := 0; i < numReaders+numWriters; i++ {
			<-done
		}

	}

}

func BenchmarkMutexSameReadWrite1(b *testing.B) {
	benchmarkMutex(b, 1, 1, 32)
}
func BenchmarkMutexSameReadWrite2(b *testing.B) {
	benchmarkMutex(b, 2, 2, 32)
}
func BenchmarkMutexSameReadWrite4(b *testing.B) {
	benchmarkMutex(b, 4, 4, 32)
}
func BenchmarkMutex1(b *testing.B) {
	benchmarkMutex(b, 2, 8, 32)
}
func BenchmarkMutex2(b *testing.B) {
	benchmarkMutex(b, 16, 4, 64)
}
func BenchmarkMutex3(b *testing.B) {
	benchmarkMutex(b, 1, 2, 128)
}
func BenchmarkMutex4(b *testing.B) {
	benchmarkMutex(b, 128, 32, 256)
}
func BenchmarkMutex5(b *testing.B) {
	benchmarkMutex(b, 1024, 2048, 64)
}
func BenchmarkMutex6(b *testing.B) {
	benchmarkMutex(b, 2048, 1024, 512)
}
