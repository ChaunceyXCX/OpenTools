package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	atotto "github.com/atotto/clipboard"
	"github.com/ChaunceyXCX/OpenTools/internal/core/launcher"
	"github.com/ChaunceyXCX/OpenTools/internal/core/scanner"
)

type jsonrpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type jsonrpcResponse struct {
	JSONRPC string `json:"jsonrpc"`
	ID      any    `json:"id"`
	Result  any    `json:"result,omitempty"`
	Error   *rpcError `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type toolDefinition struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	InputSchema any    `json:"inputSchema"`
}

type listToolsResult struct {
	Tools []toolDefinition `json:"tools"`
}

type callToolParams struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}

type callToolResult struct {
	Content []contentItem `json:"content"`
	IsError bool          `json:"isError"`
}

type contentItem struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func StartMCPServer() {
	log.Println("[MCP] server starting on stdio")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var req jsonrpcRequest
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			log.Printf("[MCP] parse error: %v", err)
			continue
		}
		resp := handleRequest(req)
		respBytes, _ := json.Marshal(resp)
		fmt.Fprintf(os.Stdout, "%s\n", string(respBytes))
		os.Stdout.Sync()
	}
}

func handleRequest(req jsonrpcRequest) jsonrpcResponse {
	switch req.Method {
	case "initialize":
		return jsonrpcResponse{
			JSONRPC: "2.0", ID: req.ID,
			Result: map[string]any{
				"protocolVersion": "2024-11-05",
				"capabilities":    map[string]any{"tools": map[string]any{}},
				"serverInfo":      map[string]any{"name": "ztools-mcp", "version": "1.0.0"},
			},
		}
	case "tools/list":
		return jsonrpcResponse{
			JSONRPC: "2.0", ID: req.ID,
			Result: listToolsResult{Tools: []toolDefinition{
				{Name: "scan_apps", Description: "Scan installed applications", InputSchema: map[string]any{"type": "object", "properties": map[string]any{}}},
				{Name: "launch_app", Description: "Launch an application by path", InputSchema: map[string]any{"type": "object", "properties": map[string]any{"path": map[string]any{"type": "string"}}, "required": []string{"path"}}},
				{Name: "get_clipboard", Description: "Get clipboard content", InputSchema: map[string]any{"type": "object", "properties": map[string]any{}}},
			}},
		}
	case "tools/call":
		var p callToolParams
		if err := json.Unmarshal(req.Params, &p); err != nil {
			return jsonrpcResponse{JSONRPC: "2.0", ID: req.ID, Error: &rpcError{Code: -32602, Message: err.Error()}}
		}
		return handleToolCall(p)
	case "notifications/initialized":
		return jsonrpcResponse{JSONRPC: "2.0", ID: req.ID, Result: map[string]any{}}
	default:
		return jsonrpcResponse{JSONRPC: "2.0", ID: req.ID, Result: map[string]any{}}
	}
}

func handleToolCall(p callToolParams) jsonrpcResponse {
	switch p.Name {
	case "scan_apps":
		cmds, err := scanner.ScanLinuxApplications()
		if err != nil {
			return jsonrpcResponse{JSONRPC: "2.0", Error: &rpcError{Code: -32603, Message: err.Error()}}
		}
		data, _ := json.Marshal(map[string]any{"apps": cmds, "count": len(cmds)})
		return jsonrpcResponse{JSONRPC: "2.0", Result: callToolResult{
			Content: []contentItem{{Type: "text", Text: string(data)}},
		}}
	case "launch_app":
		path, _ := p.Arguments["path"].(string)
		if path == "" {
			return jsonrpcResponse{JSONRPC: "2.0", Error: &rpcError{Code: -32602, Message: "path required"}}
		}
		if err := launcher.LaunchApp(path); err != nil {
			return jsonrpcResponse{JSONRPC: "2.0", Error: &rpcError{Code: -32603, Message: err.Error()}}
		}
		return jsonrpcResponse{JSONRPC: "2.0", Result: callToolResult{
			Content: []contentItem{{Type: "text", Text: `{"success":true}`}},
		}}
	case "get_clipboard":
		text, err := atotto.ReadAll()
		if err != nil {
			text = ""
		}
		data, _ := json.Marshal(map[string]string{"content": text})
		return jsonrpcResponse{JSONRPC: "2.0", Result: callToolResult{
			Content: []contentItem{{Type: "text", Text: string(data)}},
		}}
	default:
		return jsonrpcResponse{JSONRPC: "2.0", Error: &rpcError{Code: -32601, Message: "unknown tool: " + p.Name}}
	}
}

func mustJSON(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}
