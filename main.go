package main

import (
	"net/http"

	"github.com/cuvva/cuvva-public-go/lib/clog"
	"github.com/cuvva/cuvva-public-go/lib/config"
	"github.com/cuvva/cuvva-public-go/lib/crpc"
	"github.com/cuvva/cuvva-public-go/lib/middleware/request"
	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	DStore "github.com/tonymj76/cuvva-sample-code/datastore"
	Handler "github.com/tonymj76/cuvva-sample-code/handlers"
	Models "github.com/tonymj76/cuvva-sample-code/models"
)

type ServerConfig struct {
	Logging clog.Config `json:"logging"`

	Server config.Server `json:"server"`
}

func recoverMe() {
	if err := recover(); err != nil {
		logrus.Fatalln("recovering from panic ", err)
	}
}

func main() {
	defer recoverMe()
	cfg := &ServerConfig{
		Logging: clog.Config{
			Format: "text",
			Debug:  true,
		},

		Server: config.Server{
			Addr: "127.0.0.1:3003",
		},
	}
	// log := cfg.Logging.Configure()
	connect, err := DStore.NewConnection(logrus.New())
	if err != nil {
		logrus.WithError(err).Fatalln("filed to initiate database")
	}
	hdl := Handler.NewService(connect)

	// create a new RPC server
	hw := crpc.NewServer(unsafeNoAuthentication)

	// add logging middleware
	hw.Use(crpc.Logger())

	// add default instrumentation
	hw.Use(crpc.Instrument(prometheus.DefaultRegisterer))

	hw.Register("create_merchant", "2021-06-29", Models.CreateRequestSchema, hdl.CreateMerchant)

	mux := chi.NewRouter()

	mux.Use(request.RequestID)
	mux.Use(request.Logger(logrus.NewEntry(logrus.New())))

	// mount system endpoints for health and monitoring
	mux.Route("/system", func(mux chi.Router) {
		mux.Handle("/metrics", promhttp.Handler())
	})

	mux.With(request.StripPrefix("/v1")).Handle("/v1/*", hw)

	s := &http.Server{Handler: mux}

	logrus.WithField("addr", cfg.Server.Addr).Info("listening")

	if err := cfg.Server.ListenAndServe(s); err != nil {
		logrus.WithError(err).Fatal("listen failed")
	}
}

// unsafeNoAuthentication middleware that will always return next handler
func unsafeNoAuthentication(next crpc.HandlerFunc) crpc.HandlerFunc {
	return func(res http.ResponseWriter, req *crpc.Request) error {
		return next(res, req)
	}
}
