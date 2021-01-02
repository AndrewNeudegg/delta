package probes

import (
	"net/http"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// ProbeServer allows the application to state readiness and liveness independently.
type ProbeServer struct {
	isReady    bool   // is the application ready?
	isAlive    bool   // is the application alive?
	ListenAddr string // what port will the liveness probe listen on?
}

// ReadyNow states that this application is prepared to handle requests.
func (p *ProbeServer) ReadyNow() {
	p.isReady = true
}

// AliveNow states that this application is prepared to function.
func (p *ProbeServer) AliveNow() {
	p.isAlive = true
}

// StartProbeServer is a blocking simple http server.
func (p *ProbeServer) StartProbeServer() {
	http.HandleFunc("/liveness", func(rw http.ResponseWriter, r *http.Request) {
		log.WithFields(
			log.Fields{
				"isReady": p.isReady,
				"isAlive": p.isAlive,
			},
		).Debug("probes: liveness probe check")

		if p.isAlive {
			rw.WriteHeader(200)
			return
		}
		rw.WriteHeader(500)
		return
	})

	http.HandleFunc("/readiness", func(rw http.ResponseWriter, r *http.Request) {
		log.WithFields(
			log.Fields{
				"isReady": p.isReady,
				"isAlive": p.isAlive,
			},
		).Debug("probes: readiness probe check")

		if p.isReady {
			rw.WriteHeader(200)
			return
		}
		rw.WriteHeader(500)
		return
	})

	err := http.ListenAndServe(p.ListenAddr, nil)
	if err != nil {
		log.Error(errors.Wrap(err, "failed to listen and serve probes"))
	}
}
