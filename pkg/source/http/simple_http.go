package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/source"
	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
)

// httpSinkServerResponse is the response that the server will send.
type httpSinkServerResponse struct {
	ID string `json:"id"` // ID is the response ID for this accepted event.
}

// HttpSink is a http server.
type SimpleHttpSink struct {
	source.S
	ListenAddr  string
	MaxBodySize int

	inboundCh chan<- events.Event
	server    *http.Server
}

// init this sink.
func (s *SimpleHttpSink) init(ch chan<- events.Event) error {
	s.inboundCh = ch
	s.server = &http.Server{Addr: s.ListenAddr, Handler: s}
	return nil
}

// ServeHTTP allows this to become a http server.
func (s *SimpleHttpSink) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("must http post to sink"))
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
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if len(body) > s.MaxBodySize {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(fmt.Sprintf("body too large, exceeded '%d' bytes", s.MaxBodySize)))
		return
	}

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
		return
	}

	s.inboundCh <- events.EventMsg{
		ID:      uniqueID,
		Headers: r.Header,
		URI:     r.RequestURI,
		Content: body,
	}

	rw.WriteHeader(http.StatusAccepted)
	rw.Write(responseBytes)
	return
}

// Do will init this component and start the listen and serve.
func (s SimpleHttpSink) Do(ctx context.Context, ch chan<- events.Event) error {
	err := s.init(ch)
	if err != nil {
		return err
	}

	// gracefully await server shutdown..
	go func() {
		for {
			select {
			case _ = <-ctx.Done():
				s.server.Shutdown(ctx)
			}
		}
	}()

	// do the serving.
	return s.server.ListenAndServe()
}
