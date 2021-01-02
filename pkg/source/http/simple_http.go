package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/source"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
)

// httpSinkServerResponse is the response that the server will send.

const (
	// SuccessStatus is the string success message.
	SuccessStatus = "success"
	// FailureStatus is the string failed message.
	FailureStatus = "success"
)

type httpSinkServerResponse struct {
	ID     string `json:"id"`     // ID is the response ID for this accepted event.
	Reason string `json:"reason"` // Reason is why the response happened as it did.
	Status string `json:"status"` // Status states what happened to this event.
}

// SimpleHTTPSink is a http server.
type SimpleHTTPSink struct {
	source.S
	ListenAddr  string
	MaxBodySize int

	inboundCh chan<- events.Event
	server    *http.Server
}

// ID returns a human readable identifier for this thing.
func (s SimpleHTTPSink) ID() string {
	return "source/SimpleHTTPSink"
}

// init this sink.
func (s *SimpleHTTPSink) init(ch chan<- events.Event) error {
	s.inboundCh = ch
	s.server = &http.Server{Addr: s.ListenAddr, Handler: s}
	return nil
}

// ServeHTTP allows this to become a http server.
func (s *SimpleHTTPSink) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("must http post to sink"))
		return
	}

	uniqueID := uuid.New().String()

	log.Debugf("received '%s' at '%s%s'.", uniqueID, r.Host, r.RequestURI)
	responseBytes, err := json.Marshal(httpSinkServerResponse{
		ID:     uniqueID,
		Reason: "none",
		Status: SuccessStatus,
	})

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Error(errors.Wrap(err, "failed to unmarshal http message"))
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
		log.Error(errors.Wrap(err, "failed to read http body"))
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	fail := func(err error) {
		responseBytes, _ := json.Marshal(httpSinkServerResponse{
			ID:     uniqueID,
			Reason: err.Error(),
			Status: FailureStatus,
		})

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
		} else {
			rw.WriteHeader(http.StatusBadRequest)
		}

		rw.Write(responseBytes)
		wg.Done()
	}

	complete := func() {
		rw.WriteHeader(http.StatusAccepted)
		rw.Write(responseBytes)
		wg.Done()
	}

	s.inboundCh <- events.EventMsg{
		ID:           uniqueID,
		Headers:      r.Header,
		URI:          r.RequestURI,
		Content:      body,
		FailFunc:     &fail,
		CompleteFunc: &complete,
	}

	wg.Wait()
	return
}

// SDo will init this component and start the listen and serve.
func (s SimpleHTTPSink) SDo(ctx context.Context, ch chan<- events.Event) error {
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
	log.Infof("starting sink server at '%s'", s.ListenAddr)
	return s.server.ListenAndServe()
}
