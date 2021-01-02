package telemetry

import (
	"context"
	"net/http"
	// "github.com/prometheus/client_golang/prometheus/promhttp"
)

// PrometheusServer is a telemetry server for Prometheus.
type PrometheusServer struct {
	ListenAddr string
	Route      string
	server     *http.Server
}

// Serve will start the server as configured.
func (s *PrometheusServer) Serve(ctx context.Context) error {
	// s.server = &http.Server{Addr: s.ListenAddr, Handler: promhttp.Handler()}
	s.server = &http.Server{Addr: s.ListenAddr, Handler: nil}
	return s.server.ListenAndServe()

}

// Stop the server.
func (s *PrometheusServer) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
