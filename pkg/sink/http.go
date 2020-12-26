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

type HttpSinkServerConfiguration struct {
	SinkServerConfiguration
	ListenAddr string
}

type httpSinkServer struct {
	config *HttpSinkServerConfiguration
	server *http.Server
}

type httpSinkServerResponse struct {
	Id string `json:"id"`
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
		Id: uniqueID,
	})

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	maxBodySize := 512
	if len(body) > maxBodySize {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(fmt.Sprintf("body too large, exceeded '%d' bytes", maxBodySize)))
	}

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
		return
	}

	s.config.SinkServerConfiguration.ToChan <- &SunkMessage{
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

func NewHTTPSinkServer(c *HttpSinkServerConfiguration) (SinkServer, error) {
	return &httpSinkServer{
		config: c,
	}, nil
}
