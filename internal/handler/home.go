package handler

import (
	"net/http"
	"text/template"

	"github.com/obrel/go-lib/pkg/log"
)

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("index.html")
	if err != nil {
		log.For("handler", "home").Error(err)
		w.WriteHeader(500)
		w.Write([]byte("Internal server error"))
	}

	t.Execute(w, "ws://"+r.Host+"/ws")
}
