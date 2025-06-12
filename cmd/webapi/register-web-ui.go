//go:build webui

package main

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"strings"

	"github.com/AleksK26/WASA_AleksK_2024-25/webui"
)

func registerWebUI(hdl http.Handler) (http.Handler, error) {
	distDirectory, err := fs.Sub(webui.Dist, "dist")
	if err != nil {
		return nil, fmt.Errorf("error embedding WebUI dist/ directory: %w", err)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RequestURI, "/dashboard/") {
			http.StripPrefix("/dashboard/", http.FileServer(http.FS(distDirectory))).ServeHTTP(w, r)
			return
		}
		if strings.HasPrefix(r.RequestURI, "/api/") {
			hdl.ServeHTTP(w, r)
			return
		}
		// Serve index.html for all other routes (SPA fallback)
		indexFile, err := distDirectory.Open("index.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("index.html not found"))
			return
		}
		defer indexFile.Close()
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = io.Copy(w, indexFile)
	}), nil
}
