package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type fileInfo struct {
	Name    string
	Size    int64
	ModTime string
	URL     string
}

type pageData struct {
	Title      string
	UploadDir  string
	Files      []fileInfo
	Message    string
	AccessNote string
}

var tpl = template.Must(template.New("index").Parse(`<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <title>{{.Title}}</title>
  <style>
    body { font-family: Arial, sans-serif; max-width: 900px; margin: 30px auto; padding: 0 16px; }
    h1 { margin-bottom: 8px; }
    .card { border: 1px solid #ddd; border-radius: 10px; padding: 16px; margin: 16px 0; }
    input[type=file] { margin-right: 12px; }
    table { width: 100%; border-collapse: collapse; }
    th, td { text-align: left; padding: 10px; border-bottom: 1px solid #eee; }
    .muted { color: #666; }
    .ok { color: green; }
    code { background: #f5f5f5; padding: 2px 6px; border-radius: 4px; }
    .warn { background: #fff7e6; border-left: 4px solid #f0ad4e; padding: 10px; }
  </style>
</head>
<body>
  <h1>{{.Title}}</h1>
  <p class="muted">Upload directory: <code>{{.UploadDir}}</code></p>
  <p class="muted">{{.AccessNote}}</p>
  {{if .Message}}<p class="ok">{{.Message}}</p>{{end}}

  <div class="card">
    <form action="/upload" method="post" enctype="multipart/form-data">
      <input type="file" name="file" required>
      <button type="submit">Upload</button>
    </form>
  </div>

  <div class="card">
    <h3>Uploaded files</h3>
    {{if .Files}}
    <table>
      <thead>
        <tr><th>Name</th><th>Size (bytes)</th><th>Modified</th><th>Download</th></tr>
      </thead>
      <tbody>
      {{range .Files}}
        <tr>
          <td>{{.Name}}</td>
          <td>{{.Size}}</td>
          <td>{{.ModTime}}</td>
          <td><a href="{{.URL}}" target="_blank">Open</a></td>
        </tr>
      {{end}}
      </tbody>
    </table>
    {{else}}
    <p class="muted">No files uploaded yet.</p>
    {{end}}
  </div>

  <div class="warn">
    <strong>Storage mode note:</strong> the same app can be deployed with <code>emptyDir</code>, <code>hostPath</code>, or a shared PVC. The mounted path remains <code>/data/uploads</code>; only the backing storage changes.
  </div>
</body>
</html>`))

func main() {
	uploadDir := getenv("UPLOAD_DIR", "/data/uploads")
	title := getenv("APP_TITLE", "Go File Uploader")
	accessNote := getenv("ACCESS_NOTE", "Files are stored using the Kubernetes volume mounted at /data/uploads.")
	addr := getenv("LISTEN_ADDR", ":8080")

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Fatalf("create upload dir: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		files, err := listFiles(uploadDir)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := pageData{
			Title:      title,
			UploadDir:  uploadDir,
			Files:      files,
			Message:    r.URL.Query().Get("msg"),
			AccessNote: accessNote,
		}
		if err := tpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			http.Error(w, fmt.Sprintf("parse form: %v", err), http.StatusBadRequest)
			return
		}
		f, hdr, err := r.FormFile("file")
		if err != nil {
			http.Error(w, fmt.Sprintf("read upload: %v", err), http.StatusBadRequest)
			return
		}
		defer f.Close()

		safeName := filepath.Base(hdr.Filename)
		if safeName == "." || safeName == string(filepath.Separator) || safeName == "" {
			http.Error(w, "invalid file name", http.StatusBadRequest)
			return
		}
		dstPath := filepath.Join(uploadDir, safeName)
		dst, err := os.Create(dstPath)
		if err != nil {
			http.Error(w, fmt.Sprintf("create file: %v", err), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, f); err != nil {
			http.Error(w, fmt.Sprintf("save file: %v", err), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/?msg="+safeName+" uploaded successfully", http.StatusSeeOther)
	})

	fileServer := http.FileServer(http.Dir(uploadDir))
	mux.Handle("/files/", http.StripPrefix("/files/", fileServer))

	log.Printf("listening on %s, upload dir=%s", addr, uploadDir)
	if err := http.ListenAndServe(addr, logRequests(mux)); err != nil {
		log.Fatal(err)
	}
}

func listFiles(root string) ([]fileInfo, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	var out []fileInfo
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		out = append(out, fileInfo{
			Name:    info.Name(),
			Size:    info.Size(),
			ModTime: info.ModTime().Format(time.RFC3339),
			URL:     "/files/" + strings.ReplaceAll(info.Name(), " ", "%20"),
		})
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Name < out[j].Name
	})
	return out, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}
