package main

import (
	"context"
	"fmt"
	"github.com/gogo/protobuf/types"
	"github.com/maxim-kuderko/metrics-collector/internal/initializers"
	"github.com/maxim-kuderko/metrics-collector/internal/repositories"
	"github.com/maxim-kuderko/metrics-collector/internal/service"
	"github.com/maxim-kuderko/metrics-collector/proto"
	"github.com/spf13/viper"
	_ "github.com/un000/grpc-snappy"
	"go.uber.org/atomic"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"net"
	"time"
)

func main() {
	app := fx.New(
		fx.NopLogger,
		fx.Provide(
			initializers.NewConfig,
			service.NewService,
			repositories.NewStdout,
			newServer,
		),
		fx.Invoke(grpcInit),
	)
	if err := app.Start(context.Background()); err != nil {
		panic(err)
	}
}

func grpcInit(s proto.MetricsCollectorGrpcServer, v *viper.Viper) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", v.GetString(`GRPC_SERVER_PORT`)))
	if err != nil {
		panic(err)
	}
	serv := grpc.NewServer(grpc.ReadBufferSize(100 << 20))
	proto.RegisterMetricsCollectorGrpcServer(serv, s)
	if err := serv.Serve(lis); err != nil {
		panic(err)
	}

}

type server struct {
	s *service.Service
	c *atomic.Int64
	proto.UnimplementedMetricsCollectorGrpcServer
}

func newServer(s *service.Service) proto.MetricsCollectorGrpcServer {
	srv := &server{
		s: s,
		c: atomic.NewInt64(0),
	}
	go func() {
		w := 1
		t := time.NewTicker(time.Second * time.Duration(w))
		for range t.C {
			fmt.Println(fmt.Sprintf("server %0.2fm req/sec ", float64(srv.c.Swap(0))/1000000/float64(w)))
		}
	}()
	return srv
}

var emptyRes = &types.Empty{}

func (h *server) Send(metrics proto.MetricsCollectorGrpc_SendServer) error {
	for {
		m, err := metrics.Recv()
		if err != nil {
			return err
		}
		h.c.Add(m.Values.Count)

		//h.s.Send(m)
	}
}
