package memory

import (
	"net/http"
	"runtime"
	"runtime/debug"
	"strconv"

	"github.com/thanhbinhdoan1993/practice-kuard/pkg/apiutils"

	"github.com/julienschmidt/httprouter"
)

type MemoryAPI struct {
	leaks [][]byte
}

// MemoryStatus is returned from a GET to this API endpoint
type MemoryStatus struct {
	MemStats runtime.MemStats `json:"memStatus"`
}

func New() *MemoryAPI {
	return &MemoryAPI{}
}

func (m *MemoryAPI) AddRoutes(r *httprouter.Router, base string) {
	r.GET(base+"/api", m.APIGet)
	r.POST(base+"/api/alloc", m.APIAlloc)
	r.POST(base+"/api/clear", m.APIClear)
}

func (m *MemoryAPI) APIGet(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	resp := &MemoryStatus{}
	runtime.ReadMemStats(&resp.MemStats)
	apiutils.ServeJSON(w, resp)
}

func (m *MemoryAPI) APIAlloc(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sSize := r.URL.Query().Get("size")
	if len(sSize) == 0 {
		http.Error(w, "size not specified", http.StatusBadRequest)
		return
	}

	i, err := strconv.ParseInt(sSize, 10, 64)
	if err != nil {
		http.Error(w, "bad size param", http.StatusBadRequest)
		return
	}

	leak := make([]byte, i, i)
	for i := 0; i < len(leak); i++ {
		leak[i] = 'x'
	}
	m.leaks = append(m.leaks, leak)
}

func (m *MemoryAPI) APIClear(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	m.leaks = nil
	runtime.GC()
	debug.FreeOSMemory()
}
