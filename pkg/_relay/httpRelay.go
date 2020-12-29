package _relay

import (
	"context"
	"net/http"

	"container/list"

	"github.com/andrewneudegg/delta/pkg/events"
)

// ---------------------------- Types ----------------------------

type singularEvent events.EventMsg
type bulkEvent []events.EventMsg

// ---------------------------- Server ----------------------------

// HTTPRelayServer is a HTTP implementation of a relay server.
type HTTPRelayServer struct {
	Server
	http.Handler

	ListenAddress string
	server        *http.Server

	queue *list.List
}

// Start serving the HTTP relay server.
func (s *HTTPRelayServer) Start() error {
	s.server = &http.Server{Addr: s.ListenAddress, Handler: s}
	s.queue = list.New()
	return s.server.ListenAndServe()
}

// Stop HTTP relay server.
func (s *HTTPRelayServer) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *HTTPRelayServer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	// readCount := func(r *http.Request) (int, error) {
	// 	count := 1
	// 	vCount, ok := r.URL.Query()["count"]

	// 	if ok && len(vCount) != 0 {
	// 		i, err := strconv.Atoi(vCount[0])
	// 		if err != nil {
	// 			return 1, err
	// 		}
	// 		count = i
	// 	}

	// 	return count, nil
	// }

	switch r.RequestURI {
	case "/next":

	case "/complete":

	case "/fail":

	case "/load":

	default:
		rw.WriteHeader(http.StatusNotFound)
		return
	}
}

// ---------------------------- Client ----------------------------

// HTTPRelayClient is a HTTP implementation of a relay client.
type HTTPRelayClient struct {
	Client

	ServerAddress string
}
