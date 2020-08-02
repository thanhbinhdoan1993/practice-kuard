package debugprobe

import (
	"net/http"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
)

const maxHistory = 20

type Probe struct {
	basePath string
	mu       sync.Mutex

	lastID int

	c       ProbeConfig
	history []*ProbeHistory
}

type ProbeHistory struct {
	ID   int
	When time.Time
	Code int
}

func New() *Probe {
	return &Probe{}
}

func (p *Probe) AddRoutes(r *httprouter.Router, base string) {
	r.GET(base, p.Handle)
	r.GET(base+"/api", p.APIGet)
	r.PUT(base+"/api", p.APIPut)

	if p.basePath != "" {
		p.basePath = base
	}
}

func (p *Probe) APIGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.lockedGet(w, r)
}

func (p *Probe) lockedGet(w http.ResponseWriter, r *http.Request) {
	s := &ProbeStatus{
		ProbePath: p.basePath,
		FailNext:  p.c.FailNext,
	}
	l := len(p.history)
	s.History = make([]ProbeStatusHistory, l)
	for i, v := range p.history {
		h := &s.History[l-1-i]
		h.ID = v.ID
		h.When = htmlutils.FriendlyTime(v.When)
	}

}
