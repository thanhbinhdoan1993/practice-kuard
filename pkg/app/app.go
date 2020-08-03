package app

import (
	"html/template"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/thanhbinhdoan1993/practice-kuard/pkg/debugprobe"
	"github.com/thanhbinhdoan1993/practice-kuard/pkg/dnsapi"
	"github.com/thanhbinhdoan1993/practice-kuard/pkg/env"
	"github.com/thanhbinhdoan1993/practice-kuard/pkg/htmlutils"
	"github.com/thanhbinhdoan1993/practice-kuard/pkg/keygen"
	"github.com/thanhbinhdoan1993/practice-kuard/pkg/memory"
	memqserver "github.com/thanhbinhdoan1993/practice-kuard/pkg/memq/server"
	"github.com/thanhbinhdoan1993/practice-kuard/pkg/sitedata"
	"github.com/thanhbinhdoan1993/practice-kuard/pkg/version"

	"github.com/felixge/httpsnoop"
	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func init() {
	prometheus.MustRegister(requestDuration)
}

var requestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "request_duration_seconds",
	Help:    "Time serving HTTP request",
	Buckets: prometheus.DefBuckets,
}, []string{"method", "route", "status_code"})

func promMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := httpsnoop.CaptureMetrics(h, w, r)
		requestDuration.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(m.Code)).Observe(m.Duration.Seconds())
	})
}

func loggingMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

type pageContext struct {
	URLBase      string       `json:"urlBase"`
	Hostname     string       `json:"hostname"`
	Addrs        []string     `json:"addrs"`
	Version      string       `json:"version"`
	VersionColor template.CSS `json:"versionColor"`
	RequestDump  string       `json:"requestDump"`
	RequestProto string       `json:"requestProto"`
	RequestAddr  string       `json:"requestAddr"`
}

type App struct {
	c  Config
	tg *htmlutils.TemplateGroup

	m     *memory.MemoryAPI
	live  *debugprobe.Probe
	ready *debugprobe.Probe
	env   *env.Env
	dns   *dnsapi.DNSAPI
	kg    *keygen.KeyGen
	mq    *memqserver.Server

	r *httprouter.Router
}

func (k *App) getPageContext(r *http.Request, urlBase string) *pageContext {
	c := &pageContext{}
	c.URLBase = urlBase
	c.Hostname, _ = os.Hostname()

	addrs, _ := net.InterfaceAddrs()
	c.Addrs = []string{}
	for _, addr := range addrs {
		// Check the address type and if it is not a loopback
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				c.Addrs = append(c.Addrs, ipnet.IP.String())
			}
		}
	}

	c.Version = version.VERSION
	c.VersionColor = template.CSS(htmlutils.ColorFromString(version.VERSION))
	reqDump, _ := httputil.DumpRequest(r, false)
	c.RequestDump = strings.TrimSpace(string(reqDump))
	c.RequestProto = r.Proto
	c.RequestAddr = r.RemoteAddr

	return c
}

func (k *App) getRootHandler(urlBase string) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		k.tg.Render(w, "index.html", k.getPageContext(r, urlBase))
	})
}

// Exists reports whether the named file or directory exists.
func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func (k *App) Run() {
	r := promMiddleware(loggingMiddleware(k.r))

	// Look to see if we can find TLS certs
	certFile := filepath.Join(k.c.TLSDir, "kuard.crt")
	keyFile := filepath.Join(k.c.TLSDir, "kuard.key")

	if fileExists(certFile) && fileExists(keyFile) {
		go func() {
			log.Printf("Serving HTTPS on %v", k.c.TLSAddr)
			log.Fatal(http.ListenAndServeTLS(k.c.TLSAddr, certFile, keyFile, r))
		}()
	} else {
		log.Printf("Could not find certificates to serve TLS")
	}

	log.Printf("Serving on HTTP on %v", k.c.ServerAddr)
	log.Fatal(http.ListenAndServe(k.c.ServerAddr, r))
}

func NewApp() *App {
	k := &App{
		tg: &htmlutils.TemplateGroup{},
		r:  httprouter.New(),
	}

	// Init all of the subcomponents

	router := k.r
	k.m = memory.New()
	k.live = debugprobe.New()
	k.ready = debugprobe.New()
	k.env = env.New()
	k.dns = dnsapi.New()
	k.kg = keygen.New()
	k.mq = memqserver.NewServer()

	// Add handlers
	for _, prefix := range []string{"", "/a", "/b", "/c"} {
		rootHandler := k.getRootHandler(prefix)
		router.GET(prefix+"/", rootHandler)
		router.GET(prefix+"/-/*path", rootHandler)

		router.Handler(http.MethodGet, prefix+"/metrics", promhttp.Handler())

		// Add static files
		sitedata.AddRoutes(router, prefix+"/built")
		sitedata.AddRoutes(router, prefix+"/static")

		router.Handler(http.MethodGet, prefix+"/fs/*filepath", http.StripPrefix(prefix+"/fs", http.FileServer(http.Dir("/"))))

		k.m.AddRoutes(router, prefix+"/mem")
		k.live.AddRoutes(router, prefix+"/healthy")
		k.ready.AddRoutes(router, prefix+"/ready")
		k.env.AddRoutes(router, prefix+"/env")
		k.dns.AddRoutes(router, prefix+"/dns")
		k.kg.AddRoutes(router, prefix+"/keygen")
		k.mq.AddRoutes(router, prefix+"/memq/server")
	}
	return k
}
