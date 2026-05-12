package plugin

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/dop251/goja"
	esbuild "github.com/evanw/esbuild/pkg/api"
)

type PluginInfo struct {
	Name     string    `json:"name"`
	Title    string    `json:"title"`
	Version  string    `json:"version"`
	Main     string    `json:"main,omitempty"`
	Logo     string    `json:"logo,omitempty"`
	Features []Feature `json:"features,omitempty"`
}

type CompatRuntime struct {
	plugins map[string]*compatPlugin
}

type compatPlugin struct {
	info        PluginInfo
	mainJS      string
	exports     map[string]any
}

func NewCompatRuntime() *CompatRuntime {
	return &CompatRuntime{
		plugins: make(map[string]*compatPlugin),
	}
}

func (r *CompatRuntime) LoadPluginDir(dir string) (*PluginInfo, error) {
	manifestPath := filepath.Join(dir, "plugin.json")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("read plugin.json: %w", err)
	}
	var info PluginInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, fmt.Errorf("parse plugin.json: %w", err)
	}

	mainFile := info.Main
	if mainFile == "" {
		for _, f := range []string{"index.js", "main.js", "app.js"} {
			if _, err := os.Stat(filepath.Join(dir, f)); err == nil {
				mainFile = f
				break
			}
		}
	}
	if mainFile == "" {
		return nil, fmt.Errorf("plugin %s has no main entry", info.Name)
	}

	jsBytes, err := os.ReadFile(filepath.Join(dir, mainFile))
	if err != nil {
		return nil, fmt.Errorf("read main %s: %w", mainFile, err)
	}

	jsCode := string(jsBytes)

	if strings.Contains(jsCode, "require(") || strings.Contains(jsCode, "module.exports") || strings.Contains(jsCode, "exports.") {
		result := esbuild.Transform(jsCode, esbuild.TransformOptions{
			Format:  esbuild.FormatIIFE,
			Target:  esbuild.ES2020,
			Loader:  esbuild.LoaderJS,
		})
		if len(result.Errors) > 0 {
			return nil, fmt.Errorf("esbuild transform %s: %s", info.Name, result.Errors[0].Text)
		}
		jsCode = string(result.Code)
	}

	r.plugins[info.Name] = &compatPlugin{
		info:   info,
		mainJS: jsCode,
		exports: make(map[string]any),
	}

	log.Printf("[Plugin] loaded %s v%s (%s)", info.Name, info.Version, mainFile)
	return &info, nil
}

func (r *CompatRuntime) Execute(name, featureCode string, args map[string]any) (any, error) {
	p, ok := r.plugins[name]
	if !ok {
		return nil, fmt.Errorf("plugin %s not found", name)
	}

	vm := goja.New()

	ztools := vm.NewObject()

	db := vm.NewObject()
	db.Set("put", func(doc map[string]any) map[string]any {
		log.Printf("[Plugin:%s] db.put %v", name, doc)
		return map[string]any{"ok": true, "id": doc["_id"], "rev": "1-abc"}
	})
	db.Set("get", func(id string) any {
		log.Printf("[Plugin:%s] db.get %s", name, id)
		return nil
	})
	db.Set("remove", func(docOrId any) map[string]any {
		return map[string]any{"ok": true}
	})
	db.Set("bulkDocs", func(docs []any) []map[string]any {
		var results []map[string]any
		for _, d := range docs {
			if doc, ok := d.(map[string]any); ok {
				results = append(results, map[string]any{"ok": true, "id": doc["_id"]})
			}
		}
		return results
	})
	db.Set("allDocs", func(key string) []map[string]any {
		return nil
	})

	dbPromises := vm.NewObject()
	dbPromises.Set("put", func(doc map[string]any) *goja.Promise {
		prom, resolve, _ := vm.NewPromise()
		resolve(map[string]any{"ok": true, "id": doc["_id"]})
		return prom
	})
	dbPromises.Set("get", func(id string) *goja.Promise {
		prom, resolve, _ := vm.NewPromise()
		resolve(nil)
		return prom
	})
	dbPromises.Set("remove", func(docOrId any) *goja.Promise {
		prom, resolve, _ := vm.NewPromise()
		resolve(map[string]any{"ok": true})
		return prom
	})
	db.Set("promises", dbPromises)
	ztools.Set("db", db)

	ztools.Set("getAppName", func() string { return "OpenTools" })
	ztools.Set("getAppVersion", func() string { return "2.4.1" })
	ztools.Set("getNativeId", func() string { return "ztools-wails-" + name })
	ztools.Set("isMacOs", func() bool { return false })
	ztools.Set("isMacOS", func() bool { return false })
	ztools.Set("isWindows", func() bool { return false })
	ztools.Set("isLinux", func() bool { return true })
	ztools.Set("isDarkColors", func() bool { return true })
	ztools.Set("getThemeInfo", func() map[string]any {
		return map[string]any{"isDark": true, "primaryColor": "#4a90d9"}
	})
	ztools.Set("onPluginEnter", func(call goja.Callable) {
		log.Printf("[Plugin:%s] onPluginEnter registered", name)
	})
	ztools.Set("onPluginReady", func(call goja.Callable) {
		log.Printf("[Plugin:%s] onPluginReady registered", name)
	})
	ztools.Set("onPluginOut", func(call goja.Callable) {})
	ztools.Set("showNotification", func(body string) {})
	ztools.Set("showToast", func(msg string, opts map[string]any) {})

	hideFn := func() {}
	ztools.Set("hideWindow", hideFn)

	ztools.Set("setExpendHeight", func(h int) {})
	ztools.Set("setSubInput", func(cb goja.Callable, placeholder string, focus bool) {})
	ztools.Set("removeSubInput", func() {})
	ztools.Set("setSubInputValue", func(text string) {})
	ztools.Set("subInputFocus", func() {})
	ztools.Set("subInputBlur", func() {})
	ztools.Set("subInputSelect", func() {})

	clip := vm.NewObject()
	clip.Set("readText", func() string { return "" })
	clip.Set("writeText", func(text string) {})
	clip.Set("readImage", func() string { return "" })
	clip.Set("writeImage", func(img string) {})
	ztools.Set("clipboard", clip)
	ztools.Set("readClipboard", func() string { return "" })
	ztools.Set("writeClipboard", func(text string) {})

	ztools.Set("simulateKeyboardTap", func(key string, mods ...string) bool { return true })
	ztools.Set("simulateMouseMove", func(x, y int) bool { return true })
	ztools.Set("simulateMouseClick", func(x, y int) bool { return true })
	ztools.Set("simulateMouseDoubleClick", func(x, y int) bool { return true })
	ztools.Set("simulateMouseRightClick", func(x, y int) bool { return true })
	ztools.Set("sendInputEvent", func(evt map[string]any) {})

	shell := vm.NewObject()
	shell.Set("execute", func(cmd string, args []string) map[string]any {
		return map[string]any{"code": 0, "stdout": "", "stderr": ""}
	})
	ztools.Set("shell", shell)

	feature := vm.NewObject()
	feature.Set("setFeature", func(featureCode string, args map[string]any) {})
	ztools.Set("feature", feature)

	ztools.Set("getPathForFile", func(file string) string { return file })
	ztools.Set("getWindowType", func() string { return "main" })

	ztools.Set("registerTool", func(name string, handler goja.Callable) {
		log.Printf("[Plugin:%s] registered MCP tool: %s", name, name)
	})

	ztools.Set("onMainPush", func(cb map[string]any) {})

	vm.Set("window", vm.NewObject())
	vm.Get("window").ToObject(vm).Set("ztools", ztools)
	vm.Set("ztools", ztools)
	vm.Set("console", map[string]any{
		"log":   func(args ...any) { log.Printf("[Plugin:%s] %v", name, args) },
		"error": func(args ...any) { log.Printf("[Plugin:%s] ERROR: %v", name, args) },
		"warn":  func(args ...any) { log.Printf("[Plugin:%s] WARN: %v", name, args) },
	})

	windowExports := vm.NewObject()
	vm.Set("exports", windowExports)
	vm.Set("module", map[string]any{"exports": windowExports})

	_, err := vm.RunString(p.mainJS)
	if err != nil {
		return nil, fmt.Errorf("exec plugin %s: %w", name, err)
	}

	if featureCode != "" {
		exportsVal := vm.Get("window").ToObject(vm).Get("exports")
		if exportsVal != nil {
			exports := exportsVal.ToObject(vm)
			if feat := exports.Get(featureCode); feat != nil {
				featObj := feat.ToObject(vm)
				if argsVal := featObj.Get("args"); argsVal != nil {
					argsObj := argsVal.ToObject(vm)
					if enterFn := argsObj.Get("enter"); enterFn != nil {
						if fn, ok := goja.AssertFunction(enterFn); ok {
							result, err := fn(goja.Undefined(), vm.ToValue(args))
							if err == nil {
								return result.Export(), nil
							}
							return nil, fmt.Errorf("plugin %s enter error: %w", name, err)
						}
					}
				}
			}
		}
	}

	return nil, nil
}

// GetPlugin returns loaded plugin info
func (r *CompatRuntime) GetPlugin(name string) *PluginInfo {
	if p, ok := r.plugins[name]; ok {
		return &p.info
	}
	return nil
}

// ListPlugins returns all loaded plugin names
func (r *CompatRuntime) ListPlugins() []PluginInfo {
	var result []PluginInfo
	for _, p := range r.plugins {
		result = append(result, p.info)
	}
	return result
}
