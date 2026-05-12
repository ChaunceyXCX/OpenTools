package main

import (
	"fmt"
	"github.com/ChaunceyXCX/OpenTools/internal/plugin"
)

func main() {
	rt := plugin.NewJSRuntime()

	err := rt.Register("hello", `
function enter(args) {
	ztools.log("hello from JS plugin!");
	ztools.notification.show("Greeting", "Hello " + (args.name || "World"));
	return { result: "ok", name: args.name };
}
`)
	if err != nil {
		panic(err)
	}

	result, err := rt.Execute("hello", map[string]any{"name": "ZTools"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Result: %+v\n", result)
	fmt.Println("PASS")
}
