package serve

import (
	"net/http"
)

type Serve struct {
	addr string
	mux  *http.ServeMux
}

func New(addr string, mux *http.ServeMux) *Serve {
	return &Serve{
		addr: addr,
		mux:  mux,
	}
}

func (s *Serve) Start() {
	http.ListenAndServe(s.addr, s.mux)
}
