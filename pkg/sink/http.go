package sink

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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

	flowCounter            prometheus.Counter
	flowCounterClientError prometheus.Counter
	flowCounterServerError prometheus.Counter
}

type httpSinkServerResponse struct {
	ID string `json:"id"`
}

func (s *httpSinkServer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("must http post to sink"))
		s.flowCounterClientError.Inc()
		return
	}

	uniqueID := uuid.New().String()

	log.Debugf("received '%s' at '%s%s'.", uniqueID, r.Host, r.RequestURI)
	responseBytes, err := json.Marshal(httpSinkServerResponse{
		ID: uniqueID,
	})

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
		s.flowCounterServerError.Inc()
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if int64(len(body)) > s.config.MaxBodySize {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(fmt.Sprintf("body too large, exceeded '%d' bytes", s.config.MaxBodySize)))
		s.flowCounterClientError.Inc()
		return
	}

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
		s.flowCounterServerError.Inc()
		return
	}

	s.config.ServerConfiguration.ToChan <- events.EventMsg{
		ID:      uniqueID,
		Headers: r.Header,
		URI:     r.RequestURI,
		Content: body,
	}

	rw.WriteHeader(http.StatusAccepted)
	rw.Write(responseBytes)
	s.flowCounter.Inc()
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

	flowCounter := promauto.NewCounter(prometheus.CounterOpts{
		Name: "sink_msgs_received",
	})

	flowCounterClientError := promauto.NewCounter(prometheus.CounterOpts{
		Name: "sink_msgs_client_error",
	})

	flowCounterServerError := promauto.NewCounter(prometheus.CounterOpts{
		Name: "sink_msgs_server_error",
	})

	return &httpSinkServer{
		config:                 c,
		flowCounter:            flowCounter,
		flowCounterClientError: flowCounterClientError,
		flowCounterServerError: flowCounterServerError,
	}, nil
}
