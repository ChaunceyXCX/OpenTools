//go:build !darwin && !windows

package colorpicker

import (
	"fmt"
	"os/exec"
	"strings"
)

type ColorResult struct {
	Hex     string `json:"hex"`
	RGB     []int  `json:"rgb"`
	Success bool   `json:"success"`
}

func Pick() (*ColorResult, error) {
	// Linux fallback: use `colorpicker` CLI if available, or return a clear error.
	if _, err := exec.LookPath("colorpicker"); err == nil {
		out, err := exec.Command("colorpicker", "--one-shot", "--short").Output()
		if err == nil {
			hex := strings.TrimSpace(string(out))
			r, g, b := hexToRGB(hex)
			return &ColorResult{Hex: "#" + hex, RGB: []int{r, g, b}, Success: true}, nil
		}
	}

	return nil, fmt.Errorf("color picker requires a GUI environment with either: " +
		"macOS 10.15+, Windows with WebView2, or the 'colorpicker' Linux CLI tool")
}

func hexToRGB(hex string) (int, int, int) {
	if len(hex) == 6 {
		var r, g, b int
		fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
		return r, g, b
	}
	return 0, 0, 0
}
