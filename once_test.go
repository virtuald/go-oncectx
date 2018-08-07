// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package oncectx_test

import (
	"context"
	. "github.com/virtuald/go-oncectx"
	"testing"
)

type one int

func (o *one) Increment() {
	*o++
}

func run(t *testing.T, ctx context.Context, once *Once, o *one, c chan bool) {
	once.Do(ctx, func() { o.Increment() })
	if v := *o; v != 1 {
		t.Errorf("once failed inside run: %d is not 1", v)
	}
	c <- true
}

func TestOnce(t *testing.T) {
	o := new(one)
	once := new(Once)
	c := make(chan bool)
	const N = 10
	for i := 0; i < N; i++ {
		go run(t, context.Background(), once, o, c)
	}
	for i := 0; i < N; i++ {
		<-c
	}
	if *o != 1 {
		t.Errorf("once failed outside run: %d is not 1", *o)
	}
}

func TestOncePanic(t *testing.T) {
	var once Once
	func() {

		once.Do(context.Background(), func() {
			// TODO: this isn't really a good test
			defer func() {
				if r := recover(); r == nil {
					t.Fatalf("Once.Do did not panic")
				}
			}()
			panic("failed")
		})
	}()

	once.Do(context.Background(), func() {
		t.Fatalf("Once.Do called twice")
	})
}

func BenchmarkOnce(b *testing.B) {
	var once Once
	f := func() {}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			once.Do(context.Background(), f)
		}
	})
}
