package server

import (
	"encoding/json"
	allocator "github.com/lrx0014/ScalableFlake/pkg/machine"
	"net/http"
)

type HTTPServer struct {
	allocator allocator.Allocator
}

func NewHTTPServer(a allocator.Allocator) *HTTPServer {
	return &HTTPServer{allocator: a}
}

func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		http.Error(w, "tenant_id missing", http.StatusBadRequest)
		return
	}

	id, err := s.allocator.Acquire(r.Context(), tenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"uid": id,
	}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		return
	}

	return
}
