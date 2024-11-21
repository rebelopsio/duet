package lua

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

type Engine struct {
	state *lua.LState
}

func NewEngine() *Engine {
	return &Engine{
		state: lua.NewState(),
	}
}

func (e *Engine) Close() {
	if e.state != nil {
		e.state.Close()
	}
}

func (e *Engine) LoadFile(filename string) error {
	return e.state.DoFile(filename)
}

func (e *Engine) CallFunction(name string, args ...lua.LValue) (lua.LValue, error) {
	fn := e.state.GetGlobal(name)
	if fn == lua.LNil {
		return nil, fmt.Errorf("function %s not found", name)
	}

	err := e.state.CallByParam(lua.P{
		Fn:      fn,
		NRet:    1,
		Protect: true,
	}, args...)
	if err != nil {
		return nil, fmt.Errorf("error calling function %s: %w", name, err)
	}

	ret := e.state.Get(-1)
	e.state.Pop(1)
	return ret, nil
}
