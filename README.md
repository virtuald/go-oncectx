go-oncectx
==========

This is a version of golang's `sync.Once` that allows calling `Do` with a
context.

It is not expected that golang's stdlib will implement something like this, 
see https://github.com/golang/go/issues/25312

Usage is like the `sync.Once` package, except that you must pass `Do` a context.
The function called by the Do function is run in a separate goroutine and
will always run to completion, regardless of whether the context's Done
channel is closed. However, the Do function itself will return either when
f returns or when the context Done channel is closed.


Why is this needed?
-------------------

A naive approach would do something like so:

```

func (o *Object) initSomething(ctx context.Context) {
    o.once.Do(func() {
      CallSomethingWithContext(ctx)  
    })
}

```

However, a careful reading will see that if you have multiple goroutines which
call `initSomething`, a cancellation of one of the caller's contexts will cause
initialization to fail for all callers, which is probably not an intuitive 
result.
