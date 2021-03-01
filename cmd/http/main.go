package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/golang/snappy"
	jsoniter "github.com/json-iterator/go"
	"github.com/julienschmidt/httprouter"
	"github.com/maxim-kuderko/metrics-collector/internal/initializers"
	"github.com/maxim-kuderko/metrics-collector/internal/repositories"
	"github.com/maxim-kuderko/metrics-collector/internal/service"
	"github.com/maxim-kuderko/metrics-collector/pkg/requests"
	metricsEnt "github.com/maxim-kuderko/metrics/entities"
	"github.com/opentracing/opentracing-go/log"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"go.uber.org/fx"
	"io"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	go func() {
		fmt.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	app := fx.New(
		fx.NopLogger,
		fx.Provide(
			initializers.NewConfig,
			service.NewService,
			repositories.NewStdout,
			newHandler,
			route,
		),
		fx.Invoke(webserver),
	)

	if err := app.Start(context.Background()); err != nil {
		panic(err)
	}
}

func route(h *handler) *httprouter.Router {
	router := httprouter.New()
	router.GET("/health", h.Health)
	router.POST("/send", h.Send)
	return router
}

func webserver(r *httprouter.Router, v *viper.Viper) {
	log.Error(http.ListenAndServe(fmt.Sprintf(":%s", v.GetString(`HTTP_SERVER_PORT`)), r))
}

type handler struct {
	s *service.Service
}

func newHandler(s *service.Service) *handler {
	return &handler{
		s: s,
	}
}
func (h *handler) Health(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func (h *handler) Send(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()
	rr := bufio.NewReader(snappy.NewReader(r.Body))
	ok := true
	for {
		b, err := rr.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			break
		}
		m := requests.MetricPool.Get().(*metricsEnt.AggregatedMetric)
		err = jsoniter.ConfigFastest.Unmarshal(b, &m)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			jsoniter.ConfigFastest.NewEncoder(w).Encode(err)
			ok = false
			break
		}
		h.s.Send(m)
	}
	if ok {
		w.WriteHeader(fasthttp.StatusNoContent)
	}
	return
}
