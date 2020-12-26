package sink

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
)

// HTTPSinkServerConfiguration represents the configuration for a http sink
type HTTPSinkServerConfiguration struct {
	ServerConfiguration
	ListenAddr  string
	MaxBodySize int64
}

type httpSinkServer struct {
	config *HTTPSinkServerConfiguration
	server *http.Server
}

type httpSinkServerResponse struct {
	ID string `json:"id"`
}

func (s *httpSinkServer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("must http post to sink"))
		return
	}

	userAgent := r.Header["User-Agent"][0]
	contentType := r.Header["Content-Type"][0]
	uniqueID := uuid.New().String()

	log.Debugf("received '%s' at '%s%s'.", uniqueID, r.Host, r.RequestURI)
	responseBytes, err := json.Marshal(httpSinkServerResponse{
		ID: uniqueID,
	})

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if int64(len(body)) > s.config.MaxBodySize {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(fmt.Sprintf("body too large, exceeded '%d' bytes", s.config.MaxBodySize)))
	}

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
		return
	}

	s.config.ServerConfiguration.ToChan <- &SunkMessage{
		MessageID:   &uniqueID,
		Host:        &r.Host,
		ContentType: &contentType,
		UserAgent:   &userAgent,
		URI:         &r.RequestURI,
		Content:     &body,
	}

	rw.WriteHeader(http.StatusAccepted)
	rw.Write(responseBytes)
	return
}

func (s *httpSinkServer) Serve(ctx context.Context) error {
	s.server = &http.Server{Addr: s.config.ListenAddr, Handler: s}
	return s.server.ListenAndServe()
}

func (s *httpSinkServer) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// NewHTTPSinkServer is a factory method for http sink servers
func NewHTTPSinkServer(c *HTTPSinkServerConfiguration) (Server, error) {
	return &httpSinkServer{
		config: c,
	}, nil
}
