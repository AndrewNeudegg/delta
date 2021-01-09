package http1

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/andrewneudegg/delta/pkg/events"
	"github.com/andrewneudegg/delta/pkg/resources/definitions"
	"github.com/google/uuid"

	log "github.com/sirupsen/logrus"
)

// Input is simple noop.
type Input struct {
	ListenAddress string `mapstructure:"listenAddress"` // ListenAddress specifies to what the server should listen (:8080).
	MaxBodySize   int    `mapstructure:"maxBodySize"`   // MaxBodySize specifies
	inboundCh     chan<- events.Collection
	server        *http.Server
}

// ID defines what this thing is.
func (i Input) ID() string {
	return ID
}

// Type defines what type of resource this is.
func (i Input) Type() definitions.ResourceType {
	return definitions.InputType
}

func (i Input) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	r.Close = true

	if r.Method != http.MethodPost {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("must http post to sink"))
		return
	}

	uniqueID := uuid.New().String()

	log.Debugf("received '%s' at '%s%s'.", uniqueID, r.Host, r.RequestURI)

	body, _ := ioutil.ReadAll(r.Body)
	if len(body) > i.MaxBodySize {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(fmt.Sprintf("body too large, exceeded '%d' bytes", i.MaxBodySize)))
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	e := events.EventMsg{
		ID:      uniqueID,
		Headers: r.Header,
		URI:     r.RequestURI,
		Content: body,
	}

	i.inboundCh <- []events.Event{e}

	err := e.Await()
	status := "Success"
	reason := ""
	if err != nil {
		status = "Failure"
		reason = err.Error()
	}

	rw.WriteHeader(http.StatusInternalServerError)
	b, _ := json.Marshal(httpSinkServerResponse{
		ID:     uniqueID,
		Reason: reason,
		Status: status,
	})
	rw.Write(b)
}

func (i *Input) init(ch chan<- events.Collection) error {
	i.inboundCh = ch
	i.server = &http.Server{Addr: i.ListenAddress, Handler: i}
	return nil
}

// DoInput will accept collections of events, passing them into the channel.
func (i *Input) DoInput(ctx context.Context, ch chan<- events.Collection) error {
	err := i.init(ch)
	if err != nil {
		return err
	}

	// gracefully await server shutdown..
	go func() {
		for {
			select {
			case _ = <-ctx.Done():
				i.server.Shutdown(ctx)
			}
		}
	}()

	// do the serving.
	log.Infof("starting sink server at '%s'", i.ListenAddress)
	return i.server.ListenAndServe()
}
