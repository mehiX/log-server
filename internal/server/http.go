package server

import (
	"encoding/json"
	"net/http"
)

func NewHTTPServer(addr string, r http.Handler) *http.Server {

	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

func RequestsHandler() http.Handler {

	httpsrvr := newHTTPServer()

	r := http.NewServeMux()

	r.HandleFunc("/", httpsrvr.handleAll)

	return r
}

type httpServer struct {
	Log *Log
}

func (s *httpServer) handleAll(wr http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		s.handleConsume(wr, r)
	case http.MethodPost:
		s.handleProduce(wr, r)
	default:
		wr.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *httpServer) handleConsume(wr http.ResponseWriter, r *http.Request) {

	var d ConsumeRequest
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		http.Error(wr, err.Error(), http.StatusBadRequest)
		return
	}

	l, err := s.Log.Read(d.Offset)
	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := ConsumeResponse{
		Record: l,
	}

	if err := json.NewEncoder(wr).Encode(resp); err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (s *httpServer) handleProduce(wr http.ResponseWriter, r *http.Request) {

	var d ProduceRequest
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		http.Error(wr, err.Error(), http.StatusBadRequest)
		return
	}

	o, err := s.Log.Append(d.Record)
	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := ProduceResponse{
		Offset: o,
	}

	if err := json.NewEncoder(wr).Encode(resp); err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return
	}

}

func newHTTPServer() *httpServer {
	return &httpServer{
		Log: NewLog(),
	}
}

type ProduceRequest struct {
	Record Record `json:"record"`
}

type ProduceResponse struct {
	Offset uint64 `json:"offset"`
}

type ConsumeRequest struct {
	Offset uint64 `json:"offset"`
}

type ConsumeResponse struct {
	Record Record `json:"record"`
}
