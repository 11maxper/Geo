package router

import (
	"net/http/pprof"
	_ "project/docs"
	"project/internal/infrastructure/jwt"
	"project/internal/infrastructure/metrics"
	mymiddleware "project/internal/infrastructure/my_middleware"
	authController "project/internal/modules/auth/controller"
	geoController "project/internal/modules/geo/controller"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewRouter(auth authController.Auther, geo geoController.GeoProvider, jwt *jwt.Jwt) *chi.Mux {
	r := chi.NewRouter()

	// Prometheus middleware
	// Должен быть первым, чтобы измерять все запросы
	r.Use(metrics.HTTPMetricsMiddleware)

	// // Proxy
	// rp := mymiddleware.NewReverseProxy("hugo", "1313")
	// r.Use(rp.ReverseProxy)

	//Swagger
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// Prometheus /metrics
	r.Get("/metrics", promhttp.Handler().ServeHTTP)

	// Pprof
	r.Route("/mycustompath/pprof", func(pp chi.Router) {
		pp.Use(jwtauth.Verifier(jwt.Token))
		pp.Use(mymiddleware.CustomAuthenticator(jwt.Token))

		// индексы и дополнительные эндпоинты
		pp.Get("/", pprof.Index)
		pp.Get("/cmdline", pprof.Cmdline)
		pp.Get("/profile", pprof.Profile)
		pp.Get("/symbol", pprof.Symbol)
		pp.Post("/symbol", pprof.Symbol)
		pp.Get("/trace", pprof.Trace)

		// профили
		pp.Get("/allocs", pprof.Handler("allocs").ServeHTTP)
		pp.Get("/block", pprof.Handler("block").ServeHTTP)
		pp.Get("/goroutine", pprof.Handler("goroutine").ServeHTTP)
		pp.Get("/heap", pprof.Handler("heap").ServeHTTP)
		pp.Get("/mutex", pprof.Handler("mutex").ServeHTTP)
		pp.Get("/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
	})

	//Auth
	r.Post("/api/register", auth.Register)
	r.Post("/api/login", auth.Login)
	r.Group(func(protected chi.Router) {
		protected.Use(jwtauth.Verifier(jwt.Token))
		protected.Use(mymiddleware.CustomAuthenticator(jwt.Token))
		//Geo
		protected.Route("/api/address", func(r chi.Router) {
			r.Post("/search", geo.AddressSearch)
			r.Post("/geocode", geo.GeocodeSearch)
		})
	})
	return r
}
