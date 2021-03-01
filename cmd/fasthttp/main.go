package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/fasthttp/router"
	jsoniter "github.com/json-iterator/go"
	"github.com/maxim-kuderko/metrics-collector/internal/initializers"
	"github.com/maxim-kuderko/metrics-collector/internal/repositories"
	"github.com/maxim-kuderko/metrics-collector/internal/service"
	metricsEnt "github.com/maxim-kuderko/metrics/entities"
	"github.com/opentracing/opentracing-go/log"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"go.uber.org/fx"
	"io"
	"runtime"
	"sync"
)

func main() {
	app := fx.New(
		fx.NopLogger,
		fx.Provide(
			initializers.NewConfig,
			initializers.NewMetrics,
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

func route(h *handler) *router.Router {
	router := router.New()
	router.GET("/health", h.Health)
	router.POST("/send", h.Send)
	return router
}

func webserver(r *router.Router, v *viper.Viper) {
	server := fasthttp.Server{
		Handler:           r.Handler,
		TCPKeepalive:      true,
		StreamRequestBody: true,
	}
	log.Error(server.ListenAndServe(fmt.Sprintf(`:%s`, v.GetString(`HTTP_SERVER_PORT`))))
}

type handler struct {
	s *service.Service
}

func newHandler(s *service.Service) *handler {
	return &handler{
		s: s,
	}
}
func (h *handler) Health(ctx *fasthttp.RequestCtx) {

}

func (h *handler) Send(ctx *fasthttp.RequestCtx) {
	metrics := parser(ctx)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for m := range metrics {
			h.s.Send(m)
		}
	}()
	wg.Wait()
	return
}

func parser(c *fasthttp.RequestCtx) chan metricsEnt.AggregatedMetric {
	output := make(chan metricsEnt.AggregatedMetric, runtime.GOMAXPROCS(0))
	go func() {
		defer close(output)
		r := bufio.NewReader(c.RequestBodyStream())
		ok := true
		for {
			b, err := r.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Println(err)
				break
			}
			var m metricsEnt.AggregatedMetric
			err = jsoniter.ConfigFastest.Unmarshal(b, &m)
			if err != nil {
				c.SetStatusCode(fasthttp.StatusBadRequest)
				jsoniter.ConfigFastest.NewEncoder(c).Encode(err)
				ok = false
				break
			}
			output <- m
		}
		if ok {
			c.SetStatusCode(fasthttp.StatusNoContent)
		}
	}()

	return output
}
