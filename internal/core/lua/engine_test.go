package lua

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestEngine(t *testing.T) {
	t.Run("LoadAndExecute", func(t *testing.T) {
		engine := NewEngine()
		defer engine.Close()

		script := `
			function test()
				return "hello"
			end
		`

		err := engine.state.DoString(script)
		if err != nil {
			t.Fatalf("Failed to load script: %v", err)
		}

		result, err := engine.CallFunction("test")
		if err != nil {
			t.Fatalf("Failed to call function: %v", err)
		}

		if result.String() != "hello" {
			t.Errorf("Expected 'hello', got %s", result.String())
		}
	})

	t.Run("RegisterAndCallGoFunction", func(t *testing.T) {
		engine := NewEngine()
		defer engine.Close()

		// Register a test function
		engine.state.SetGlobal("go_test", engine.state.NewFunction(func(L *lua.LState) int {
			L.Push(lua.LString("success"))
			return 1
		}))

		script := `
			function test()
				return go_test()
			end
		`

		err := engine.state.DoString(script)
		if err != nil {
			t.Fatalf("Failed to load script: %v", err)
		}

		result, err := engine.CallFunction("test")
		if err != nil {
			t.Fatalf("Failed to call function: %v", err)
		}

		if result.String() != "success" {
			t.Errorf("Expected 'success', got %s", result.String())
		}
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		engine := NewEngine()
		defer engine.Close()

		// Test invalid syntax
		err := engine.state.DoString("invalid lua syntax }")
		if err == nil {
			t.Error("Expected error for invalid syntax, got nil")
		}

		// Test calling non-existent function
		_, err = engine.CallFunction("nonexistent")
		if err == nil {
			t.Error("Expected error for non-existent function, got nil")
		}
	})
}
