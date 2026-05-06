package lua

import (
	"sync"

	"github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

// Pool manages a pool of Lua VMs
type Pool struct {
	pool   sync.Pool
	bridge *Bridge
	logger *zap.Logger
}

// NewPool creates a new Lua VM pool
func NewPool(logger *zap.Logger, poolSize int) *Pool {
	bridge := NewBridge(logger)

	p := &Pool{
		bridge: bridge,
		logger: logger.Named("lua-pool"),
	}

	p.pool.New = func() interface{} {
		L := lua.NewState(lua.Options{
			IncludeGoStackTrace: true,
		})
		bridge.SetupSandbox(L)
		return L
	}

	return p
}

// Get retrieves a Lua VM from the pool
func (p *Pool) Get() *lua.LState {
	return p.pool.Get().(*lua.LState)
}

// Put returns a Lua VM to the pool
func (p *Pool) Put(L *lua.LState) {
	// Clean up the VM state before returning to pool
	L.SetTop(0)
	p.pool.Put(L)
}

// Close closes all VMs in the pool
func (p *Pool) Close() {
	// Note: sync.Pool doesn't provide a way to close all instances
	// VMs will be garbage collected when no longer referenced
}