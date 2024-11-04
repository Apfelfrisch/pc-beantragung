package server

import (
	"net/http"
	"pc-beantragung/cmd/web"

	"github.com/a-h/templ"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /", templ.Handler(web.Base()))
	mux.HandleFunc("GET /signon", web.ListSignonsHandler(s.db))
	mux.HandleFunc("PUT /signon/{id}", web.UpdateSignonHandler(s.db))
	mux.HandleFunc("GET /signon/sidebar/{id}", web.ToggleSidebarHandler(s.db, true))
	mux.HandleFunc("DELETE /signon/sidebar/{id}", web.ToggleSidebarHandler(s.db, false))
	mux.HandleFunc("POST /signon/import-file", web.UploadFileHandler(s.db))

	fileServer := http.FileServer(http.FS(web.Files))
	mux.Handle("GET /assets/", fileServer)

	return mux
}
