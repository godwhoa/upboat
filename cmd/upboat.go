package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/alexedwards/scs"
	"github.com/go-chi/chi"
	"github.com/godwhoa/upboat/pkg/api"
	"github.com/godwhoa/upboat/pkg/api/middleware"
	"github.com/godwhoa/upboat/pkg/postgres"
	"github.com/godwhoa/upboat/pkg/posts"
	"github.com/godwhoa/upboat/pkg/users"
	openzipkin "github.com/openzipkin/zipkin-go"
	zhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"go.opencensus.io/exporter/zipkin"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
)

func key() string {
	k := make([]byte, 64)
	rand.Read(k)
	return base64.StdEncoding.EncodeToString(k)
}

func main() {
	localEndpoint, _ := openzipkin.NewEndpoint("upboat", "192.168.1.5:5454")

	reporter := zhttp.NewReporter("http://localhost:9411/api/v2/spans")
	defer reporter.Close()

	exporter := zipkin.NewExporter(reporter, localEndpoint)
	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	view.SetReportingPeriod(1 * time.Second)

	// setup logger
	cfg := zap.NewProductionConfig()
	// cfg.Encoding = "console"
	log, _ := cfg.Build()
	defer log.Sync()

	// setup platform dependencies
	sessionManager := scs.NewCookieManager(key())
	repos, err := postgres.NewFromOptions(postgres.Options{
		Host:   "localhost",
		DBName: "upboat",
		Port:   5432,
		User:   "postgres",
		Pass:   "bingbong",
	})
	if err != nil {
		log.Fatal("postgres.NewFromOptions", zap.Error(err))
	}
	// setup services
	us := users.NewService(repos.UserRepo)
	us = users.Chain(us, users.Logging(log), users.Tracing)
	ps := posts.NewService(repos.PostRepo)
	ps = posts.Chain(ps, posts.Logging(log), posts.Tracing)
	usersapi := api.NewUsersAPI(us, sessionManager, log)
	postsapi := api.NewPostsAPI(ps, log)
	// setup handlers
	r := chi.NewRouter()
	r.Route("/v1/api/", func(r chi.Router) {
		r.Use(sessionManager.Use)
		r.Route("/users", func(r chi.Router) {
			r.Post("/", usersapi.Register)
			r.Post("/login", usersapi.Login)
			r.Post("/logout", usersapi.Logout)
		})
		r.Route("/posts", func(r chi.Router) {
			r.Use(middleware.Auth(sessionManager))
			// CRUD posts
			r.Post("/", postsapi.Create)
			r.Group(func(r chi.Router) {
				r.Use(middleware.PostID)
				r.Get("/{postID}", postsapi.Get)
				r.Put("/{postID}", postsapi.Update)
				r.Delete("/{postID}", postsapi.Delete)
				// CRUD vote
				r.Get("/{postID}/vote", postsapi.Votes)
				r.Post("/{postID}/vote", postsapi.Vote)
				r.Delete("/{postID}/vote", postsapi.Unvote)
			})
		})
	})
	r.Get("/v1/map", func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(maproutes(r))
	})
	log.Info("Started!")
	err = http.ListenAndServe(":8080", &ochttp.Handler{Handler: r})
	log.Fatal("ListenAndServe", zap.Error(err))
}
