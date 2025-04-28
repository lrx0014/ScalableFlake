package server

import (
	"encoding/json"
	allocator "github.com/lrx0014/ScalableFlake/pkg/machine"
	"github.com/lrx0014/ScalableFlake/pkg/snowflake"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type HTTPServer struct {
	allocator allocator.Allocator
}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{}
}

func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uid, err := snowflake.GenerateUID()
	if err != nil {
		log.Errorf("failed to generate uid: %v", err)
		return
	}

	resp := map[string]interface{}{
		"uid": uid,
	}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		return
	}

	return
}
