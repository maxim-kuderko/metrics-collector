package main

import (
	"bytes"
	"context"
	"fmt"
	marshaler "github.com/golang/protobuf/proto"
	"github.com/klauspost/compress/snappy"
	"github.com/maxim-kuderko/metrics-collector/internal/initializers"
	"github.com/maxim-kuderko/metrics-collector/internal/repositories"
	"github.com/maxim-kuderko/metrics-collector/internal/service"
	"github.com/maxim-kuderko/metrics-collector/proto"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"io/ioutil"
	"net"
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
			repositories.NewInfluxDB,
		),
		fx.Invoke(webserver),
	)

	if err := app.Start(context.Background()); err != nil {
		panic(err)
	}
}
func webserver(s *service.Service, v *viper.Viper) {
	addrS := net.UDPAddr{
		Port: v.GetInt(`UDP_SERVER_PORT`),
		IP:   net.ParseIP("0.0.0.0"),
	}
	ln, _ := net.ListenUDP(`udp`, &addrS)
	buffSize := v.GetInt(`UDP_BUFFER_SIZE`)
	for {
		buff := make([]byte, buffSize)
		n, err := ln.Read(buff)
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleBuffer(buff[0:n], s)
	}
}

func handleBuffer(buff []byte, s *service.Service) {
	tmp := proto.MetricRequestPool.Get().(*proto.MetricsRequest)
	defer proto.MetricRequestPool.Put(tmp)
	r := snappy.NewReader(bytes.NewBuffer(buff))
	b, _ := ioutil.ReadAll(r)
	if err := marshaler.Unmarshal(b, tmp); err != nil {
		fmt.Println(err)
		return
	}
	for _, m := range tmp.Metrics {
		s.Send(m)
	}
}
