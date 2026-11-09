package server

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rabidaudio/inshrimpentation/api"
	"github.com/rabidaudio/inshrimpentation/database"
	"github.com/rabidaudio/inshrimpentation/live"
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

	// serve static files from dist outside of /api path
	s := chi.NewRouter()
	s.Mount("/api", r)
	s.Get("/live.mp4", func(rw http.ResponseWriter, req *http.Request) {
		r, wait, err := live.Stream(req.Context())
		if err != nil {
			log.Printf("error opening stream: %v", err)
			rw.WriteHeader(500)
		} else {
			rw.Header().Set("Content-Type", "video/mp4")
			rw.WriteHeader(200)
			io.Copy(rw, r)
			wait()
		}
	})
	s.Mount("/", http.FileServer(http.Dir("frontend/dist")))
	return s
}

type healthResponseBody struct {
	Environment api.Environment `json:"environment"`
	CurrentTime time.Time       `json:"current_time"`
}

// HealthCheck implements a simple endpoint for checking if the server is running.
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	api.Respond(w, &healthResponseBody{
		Environment: api.Env(),
		CurrentTime: time.Now(),
	}, "OK")
}
