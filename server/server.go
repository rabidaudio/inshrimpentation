package server

import (
	"crypto/tls"
	"log"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/go-chi/httplog/v3"

	"github.com/go-chi/chi/v5"
	"github.com/rabidaudio/api"
	"github.com/rabidaudio/api/env"
	"github.com/rabidaudio/inshrimpentation/database"
)

type RouteBuilder func(r chi.Router)

var routers = []RouteBuilder{
	// ADD ROUTES HERE
}

func NewServer(db *database.DB) http.Handler {
	r := chi.NewRouter()

	// attach db to context
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := database.WithDB(r.Context(), db)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	// health check
	r.Get("/", HealthCheck)
	// routes from other packages
	for _, rb := range routers {
		rb(r)
	}

	s := chi.NewRouter()

	// Logging middleware
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	s.Use(httplog.RequestLogger(logger, nil))

	// reverse-proxy requests at /live/* to mediamtx
	livestream(s, "/live", "https://localhost:8892/cam")

	// api
	s.Mount("/api", r)

	// serve static files from dist outside of /api path
	s.Mount("/", http.FileServer(http.Dir("frontend/dist")))

	return s
}

type healthResponseBody struct {
	Environment env.Environment `json:"environment"`
	CurrentTime time.Time       `json:"current_time"`
}

// HealthCheck implements a simple endpoint for checking if the server is running.
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	api.Respond(w, &healthResponseBody{
		Environment: env.Env(),
		CurrentTime: time.Now(),
	}, "OK")
}

func livestream(mux *chi.Mux, path, upstreamurl string) {
	upstream, err := url.Parse(upstreamurl)
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(upstream)
	proxy.ErrorLog = log.Default()
	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // allow self-signed cert
		},
	}
	mux.Handle(path, http.RedirectHandler(path+"/", http.StatusTemporaryRedirect))
	mux.Handle(path+"/*", http.StripPrefix(path, proxy))
}
