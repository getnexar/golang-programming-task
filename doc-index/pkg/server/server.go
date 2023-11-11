package server

import (
	"net/http"
	"time"
)

type HttpServer struct {
	server *http.Server
}

func NewHttpServer(
	address string,
	handler http.Handler,
) *HttpServer {
	return &HttpServer{
		server: &http.Server{
			Addr:    address,
			Handler: handler,
		},
	}
}

func (s *HttpServer) Start() error {
	errChan := make(chan error)
	go func() {
		errChan <- s.server.ListenAndServe()
	}()

	select {
	case err := <-errChan:
		return err
	case <-time.After(3 * time.Second):
		return nil
	}
}

func (s *HttpServer) Stop() error {
	return s.server.Close()
}
