package server

import (
	"crypto/ecdsa"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	jsoniter "github.com/json-iterator/go"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary


func apiRouter(logger zerolog.Logger, pk *ecdsa.PrivateKey, rdb redis.UniversalClient) chi.Router {
	r := chi.NewRouter()

	r.Use(httprate.LimitByIP(100, time.Minute))
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	// route auth
	r.Route("/auth", func(r chi.Router) {
		r.Post("/", handleAuth(logger, pk))
		r.Post("/register", handleRegister(logger))
	})

	// private api
	r.Group(func(r chi.Router) {
		// todo use auth middleware
		// route reviews
		r.Route("/review", func(r chi.Router) {
			r.Post("/friends", handleFriendReview(logger, rdb))
			r.Post("/smurf", handleSmurfReview(logger))
		})
		//route stats
		r.Route("/stats", func(r chi.Router) {
			r.Get("/", handleStats(logger, rdb))
		})
	})

	return r
}

func frontRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/*", http.StripPrefix("/", http.FileServer(http.Dir("./api/front/build"))).ServeHTTP)
	return r
}

func NewRouter(logger zerolog.Logger, pk *ecdsa.PrivateKey, rdb redis.UniversalClient) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.SetHeader("Access-Control-Allow-Origin", "*"))

	r.Mount("/api", apiRouter(logger, pk, rdb))
	r.Mount("/", frontRouter())

	//* раскомментить для ssr (костыль, мб можно лучше)
	//* мб как можно лучше: заставить vite все html файлы в отдельную папку и тогда для ssr будут отдельные роуты
	// r.Get("/*", http.StripPrefix("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	path := r.URL.Path
	// 	if path == "" {
	// 		path = "/index.html"
	// 	}
	// 	if !strings.HasSuffix(path, ".html") {
	// 		logger.Info().Msg("not html")
	// 		http.FileServer(http.Dir("./api/front/build")).ServeHTTP(w, r)
	// 		return
	// 	}
	// 	templ, err := template.ParseFiles("./api/front/build" + path)
	// 	logger.Info().Str("path", path).Msg("")
	// 	if err != nil {
	// 		w.WriteHeader(http.StatusNotFound)
	// 		return
	// 	}
	// 	w.WriteHeader(http.StatusOK)
	// 	templ.Execute(w, "asd")
	// })).ServeHTTP)

	return r
}
