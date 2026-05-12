package httpserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

type Server struct {
	server *http.Server
	port   int
}

type statusResponse struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Running bool   `json:"running"`
}

type apiResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func New(port int) *Server {
	mux := http.NewServeMux()
	s := &Server{
		port: port,
		server: &http.Server{
			Addr:    fmt.Sprintf("127.0.0.1:%d", port),
			Handler: mux,
		},
	}

	mux.HandleFunc("/api/status", s.handleStatus)
	mux.HandleFunc("/api/commands", s.handleCommands)
	return s
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return err
	}
	port := listener.Addr().(*net.TCPAddr).Port
	s.port = port
	writePortFile(port)
	go s.server.Serve(listener)
	log.Printf("[HTTP] server started on 127.0.0.1:%d", port)
	return nil
}

func (s *Server) Stop() error {
	removePortFile()
	return s.server.Close()
}

func (s *Server) Port() int {
	return s.port
}

func writePortFile(port int) {
	dir := filepath.Join(os.TempDir(), "opentools")
	os.MkdirAll(dir, 0755)
	os.WriteFile(filepath.Join(dir, "http-port"), []byte(fmt.Sprintf("%d", port)), 0644)
}

func removePortFile() {
	os.Remove(filepath.Join(os.TempDir(), "opentools", "http-port"))
}

func jsonResp(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apiResponse{Success: true, Data: data})
}

func jsonError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(apiResponse{Success: false, Error: msg})
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	jsonResp(w, statusResponse{
		Name:    "OpenTools",
		Version: "2.4.1",
		Running: true,
	})
}

func (s *Server) handleCommands(w http.ResponseWriter, r *http.Request) {
	jsonResp(w, []string{"scan", "launch", "clipboard"})
}
