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
	"github.com/maxim-kuderko/metrics-collector/proto"
	"github.com/opentracing/opentracing-go/log"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"io"
	"net/http"
	"sync"
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
	metrics := parser(w, r)
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

func parser(w http.ResponseWriter, r *http.Request) chan *proto.Metric {
	output := make(chan *proto.Metric, 100)
	go func() {
		defer close(output)
		defer r.Body.Close()
		r := bufio.NewReader(snappy.NewReader(r.Body))
		for {
			b, err := r.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Println(err)
				break
			}
			if len(b) < 2 {
				continue
			}
			m := proto.MetricPool.Get().(*proto.Metric)
			err = jsoniter.ConfigFastest.Unmarshal(b, m)
			if err != nil {
				continue
			}
			output <- m
		}
	}()

	return output
}
