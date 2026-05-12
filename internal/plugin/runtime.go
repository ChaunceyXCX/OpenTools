package plugin

import (
	"fmt"
	"log"

	"github.com/dop251/goja"
)

type JSRuntime struct {
	vm      *goja.Runtime
	plugins map[string]*jsPlugin
}

type jsPlugin struct {
	name    string
	code    string
	enter   goja.Callable
}

func NewJSRuntime() *JSRuntime {
	return &JSRuntime{
		plugins: make(map[string]*jsPlugin),
	}
}

func (r *JSRuntime) Register(name, code string) error {
	vm := goja.New()
	r.injectAPI(vm, name)
	_, err := vm.RunString(code)
	if err != nil {
		return fmt.Errorf("js plugin %s parse error: %w", name, err)
	}

	var p = &jsPlugin{name: name, code: code}

	if enterVal := vm.Get("enter"); enterVal != nil {
		if fn, ok := goja.AssertFunction(enterVal); ok {
			p.enter = fn
		}
	}

	vm.Set("exports", map[string]any{})
	r.plugins[name] = p
	log.Printf("[JSRuntime] registered plugin: %s", name)
	return nil
}

func (r *JSRuntime) Execute(name string, args map[string]any) (any, error) {
	p, ok := r.plugins[name]
	if !ok {
		return nil, fmt.Errorf("plugin %s not found", name)
	}
	vm := goja.New()
	r.injectAPI(vm, name)
	if _, err := vm.RunString(p.code); err != nil {
		return nil, fmt.Errorf("execute plugin %s: %w", name, err)
	}
	if enterVal := vm.Get("enter"); enterVal != nil {
		if fn, ok := goja.AssertFunction(enterVal); ok {
			result, err := fn(goja.Undefined(), vm.ToValue(args))
			if err != nil {
				return nil, fmt.Errorf("plugin %s enter error: %w", name, err)
			}
			return result.Export(), nil
		}
	}
	return nil, nil
}

func (r *JSRuntime) injectAPI(vm *goja.Runtime, name string) {
	ztools := vm.NewObject()

	dbObj := vm.NewObject()
	dbObj.Set("get", func(key string) (string, error) { return "", nil })
	dbObj.Set("put", func(key, value string) error { return nil })
	ztools.Set("db", dbObj)

	clipObj := vm.NewObject()
	clipObj.Set("readText", func() (string, error) { return "", nil })
	clipObj.Set("writeText", func(text string) error { return nil })
	ztools.Set("clipboard", clipObj)

	notifyObj := vm.NewObject()
	notifyObj.Set("show", func(title, body string) {
		log.Printf("[Plugin:%s] notification: %s - %s", name, title, body)
	})
	ztools.Set("notification", notifyObj)

	ztools.Set("log", func(args ...any) {
		log.Printf("[Plugin:%s] %v", name, args)
	})

	vm.Set("ztools", ztools)
	vm.Set("console", map[string]any{
		"log":   func(args ...any) { log.Printf("[JS:%s] %v", name, args) },
		"error": func(args ...any) { log.Printf("[JS:%s] ERROR: %v", name, args) },
	})
}
