// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package oncectx

import (
	"context"
	"sync"
)

// Once is an object that will perform exactly one action.
type Once struct {
	once sync.Once
	done chan struct{}
}

// Do calls the function f if and only if Do is being called for the
// first time for this instance of Once. In other words, given
// 	var once Once
// if once.Do(f) is called multiple times, only the first call will invoke f,
// even if f has a different value in each invocation. A new instance of
// Once is required for each function to execute.
//
// Do is intended for initialization that must be run exactly once. Since f
// is niladic, it may be necessary to use a function literal to capture the
// arguments to a function to be invoked by Do:
// 	config.once.Do(ctx, func() { config.init(filename) })
//
// Because no call to Do returns until the one call to f returns, if f causes
// Do to be called, it will deadlock.
//
// The function called by the Do function is run in a separate goroutine and
// will always run to completion, regardless of whether the context's Done
// channel is closed. However, the Do function itself will return either when
// f returns or when the context Done channel is closed.
//
// For the most intuitive results, the function called by Do should not use the
// context passed to Do.
//
// If f panics, your program the same as if any goroutine panics
//
func (o *Once) Do(ctx context.Context, f func()) {
	o.once.Do(func() {
		o.done = make(chan struct{}, 1)
		go func() {
			defer func() {
				close(o.done)
			}()
			f()
		}()
	})

	select {
	case <-o.done:
	case <-ctx.Done():
	}
}
