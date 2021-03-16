package main

import (
	"context"
	"encoding/binary"
	"fmt"
	marshaler "github.com/golang/protobuf/proto"
	"github.com/maxim-kuderko/metrics-collector/internal/initializers"
	"github.com/maxim-kuderko/metrics-collector/internal/repositories"
	"github.com/maxim-kuderko/metrics-collector/internal/service"
	"github.com/maxim-kuderko/metrics-collector/proto"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"net"
	"net/http"
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
		go handleBuffer(buff[0:n], n, s)
	}
}

func handleBuffer(buff []byte, n int, s *service.Service) {
	scanned := 0
	for scanned+4 < n {
		size := int(binary.BigEndian.Uint32(buff[scanned : scanned+4]))
		scanned += 4
		tmp := proto.Metric{}
		if err := marshaler.Unmarshal(buff[scanned:scanned+size], &tmp); err != nil {
			fmt.Println(err)
			continue
		}
		scanned += size
		s.Send(&tmp)
	}
}
