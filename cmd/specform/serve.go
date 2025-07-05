package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func NewServeCommand() *cobra.Command {
	var outputDir string
	var port int
	var verbose bool

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Serve compiled prompt specs and snapshots over HTTP",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := NewLogger(verbose)
			logger.Info("Starting specform server", "dir", outputDir, "port", port)

			withCORS := func(h http.HandlerFunc) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					logger.Debug("Received request", "method", r.Method, "path", r.URL.Path)
					w.Header().Set("Access-Control-Allow-Origin", "*")
					w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
					w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
					if r.Method == http.MethodOptions {
						w.WriteHeader(http.StatusNoContent)
						return
					}
					h(w, r)
				}
			}

			http.HandleFunc("/", withCORS(func(w http.ResponseWriter, r *http.Request) {
				info := map[string]interface{}{
					"name":        "specform Server",
					"description": "Serves compiled prompt and snapshot JSON files",
					"endpoints": map[string]string{
						"/prompts":       "List all compiled prompts",
						"/prompts/:id":   "Get a compiled prompt",
						"/snapshots":     "List all snapshots",
						"/snapshots/:id": "Get a snapshot",
					},
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(info)
				logger.Debug("Serving root endpoint")
			}))

			http.HandleFunc("/healthz", withCORS(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
				logger.Debug("Health check response", "status", http.StatusOK)
			}))

			http.HandleFunc("/prompts", withCORS(func(w http.ResponseWriter, r *http.Request) {
				logger.Info("Handling request", "method", r.Method, "path", r.URL.Path)
				var files []string
				_ = filepath.Walk(outputDir, func(path string, info fs.FileInfo, err error) error {
					if !info.IsDir() && strings.HasSuffix(info.Name(), ".prompt.json") {
						rel, _ := filepath.Rel(outputDir, path)
						files = append(files, strings.TrimSuffix(rel, ".prompt.json"))
					}
					return nil
				})
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]any{
					"count":   len(files),
					"prompts": files,
				})
			}))

			http.HandleFunc("/prompts/", withCORS(func(w http.ResponseWriter, r *http.Request) {
				logger.Info("Handling request", "method", r.Method, "path", r.URL.Path)
				id := strings.TrimPrefix(r.URL.Path, "/prompts/")
				filePath := filepath.Join(outputDir, id+".prompt.json")
				data, err := os.ReadFile(filePath)
				if err != nil {
					logger.Error("Error reading prompt file", "file", filePath, "error", err)
					http.Error(w, "Prompt not found", http.StatusNotFound)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write(data)
				logger.Debug("Serving prompt", "id", id)
			}))

			http.HandleFunc("/snapshots", withCORS(func(w http.ResponseWriter, r *http.Request) {
				logger.Info("Handling request", "method", r.Method, "path", r.URL.Path)
				var files []string
				_ = filepath.Walk(outputDir, func(path string, info fs.FileInfo, err error) error {
					if !info.IsDir() && strings.HasSuffix(info.Name(), ".snap.json") {
						rel, _ := filepath.Rel(outputDir, path)
						files = append(files, strings.TrimSuffix(rel, ".snap.json"))
					}
					return nil
				})
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"count":     len(files),
					"snapshots": files,
				})
			}))

			http.HandleFunc("/snapshots/", withCORS(func(w http.ResponseWriter, r *http.Request) {
				logger.Info("Handling request", "method", r.Method, "path", r.URL.Path)
				id := strings.TrimPrefix(r.URL.Path, "/snapshots/")
				filePath := filepath.Join(outputDir, id+".snap.json")
				data, err := os.ReadFile(filePath)
				if err != nil {
					logger.Error("Error reading snapshot file", "file", filePath, "error", err)
					http.Error(w, "Snapshot not found", http.StatusNotFound)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write(data)
			}))

			logger.Info("Serving files", "dir", outputDir, "port", port)
			fmt.Printf("Serving from %s on http://localhost:%d\n", outputDir, port)
			return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
		},
	}

	cmd.Flags().StringVarP(&outputDir, "dir", "d", ".specform", "Directory to serve from")
	cmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to serve on")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	return cmd
}
