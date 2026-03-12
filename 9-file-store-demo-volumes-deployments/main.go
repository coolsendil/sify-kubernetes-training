package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const storageDir = "/data/storage"

type WriteTextRequest struct {
	FileName string `json:"file_name"`
	Content  string `json:"content"`
}

type FileInfo struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"mod_time"`
}

func main() {
	if err := os.MkdirAll(storageDir, 0o755); err != nil {
		log.Fatalf("failed to create storage dir: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/write-text", writeTextHandler)
	mux.HandleFunc("/upload", uploadHandler)
	mux.HandleFunc("/files", listFilesHandler)
	mux.HandleFunc("/files/", readFileHandler)

	addr := ":8080"
	log.Printf("server starting on %s", addr)
	log.Printf("storage path: %s", storageDir)

	server := &http.Server{
		Addr:              addr,
		Handler:           loggingMiddleware(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status":  "ok",
		"storage": storageDir,
	})
}

func writeTextHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req WriteTextRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	req.FileName = sanitizeFileName(req.FileName)
	if req.FileName == "" {
		http.Error(w, "file_name is required", http.StatusBadRequest)
		return
	}

	targetPath := filepath.Join(storageDir, req.FileName)
	if err := os.WriteFile(targetPath, []byte(req.Content), 0o644); err != nil {
		http.Error(w, fmt.Sprintf("failed to write file: %v", err), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"message":   "text file stored successfully",
		"file_name": req.FileName,
		"path":      targetPath,
		"bytes":     len(req.Content),
	})
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, fmt.Sprintf("failed to parse multipart form: %v", err), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "missing form field 'file'", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileName := sanitizeFileName(header.Filename)
	if override := sanitizeFileName(r.FormValue("file_name")); override != "" {
		fileName = override
	}
	if fileName == "" {
		http.Error(w, "invalid filename", http.StatusBadRequest)
		return
	}

	targetPath := filepath.Join(storageDir, fileName)
	dst, err := os.Create(targetPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create destination file: %v", err), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	written, err := io.Copy(dst, file)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to save uploaded file: %v", err), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"message":   "uploaded file stored successfully",
		"file_name": fileName,
		"path":      targetPath,
		"bytes":     written,
	})
}

func listFilesHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/files" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	entries, err := os.ReadDir(storageDir)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to read storage dir: %v", err), http.StatusInternalServerError)
		return
	}

	files := make([]FileInfo, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		files = append(files, FileInfo{
			Name:    info.Name(),
			Size:    info.Size(),
			ModTime: info.ModTime(),
		})
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"storage": storageDir,
		"count":   len(files),
		"files":   files,
	})
}

func readFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := strings.TrimPrefix(r.URL.Path, "/files/")
	name = sanitizeFileName(name)
	if name == "" {
		http.Error(w, "invalid file name", http.StatusBadRequest)
		return
	}

	targetPath := filepath.Join(storageDir, name)
	data, err := os.ReadFile(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "file not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("failed to read file: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func sanitizeFileName(name string) string {
	name = filepath.Base(strings.TrimSpace(name))
	if name == "." || name == "/" || name == "" {
		return ""
	}
	return name
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s from=%s duration=%s", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
	})
}
