package app

import (
	"context"
	"net/http"
)

type Server struct {
	httpsrv *http.Server
}

func NewServer() *Server {
	return &Server{}
}

func (srv *Server) RunServer(runAddress string, handler http.Handler) error {
	srv.httpsrv = &http.Server{
		Addr:    runAddress,
		Handler: handler,
	}

	return srv.httpsrv.ListenAndServe()
}

func (srv *Server) ShutdownServer(ctx context.Context) error {
	return srv.httpsrv.Shutdown(ctx)
}
