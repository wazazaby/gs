# gs
Package `gs` exposes some of `sync` package's types and subpackages with a generic layer built on top.

Improved DevX (compile-time checks, and you avoid dealing with `any`/`interface{}`).

## Supported types
* [sync.Map](https://pkg.go.dev/sync#Map)
* [sync.Pool](https://pkg.go.dev/sync#Pool)
## Supported subpackages
* [singleflight](https://pkg.go.dev/golang.org/x/sync/singleflight)
* [atomic](https://pkg.go.dev/sync/atomic) (partially for now, only for `atomic.Value`)
## Other
`atomic.CloseSafeChan` is a chan wrapper that guarantees safe concurrent closing operations on the channel.
