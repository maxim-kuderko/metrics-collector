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
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"io"
	"net"
	_ "net/http/pprof"
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
	serv := grpc.NewServer()
	proto.RegisterMetricsCollectorGrpcServer(serv, s)
	if err := serv.Serve(lis); err != nil {
		panic(err)
	}

}

type server struct {
	s *service.Service
	proto.UnimplementedMetricsCollectorGrpcServer
}

func newServer(s *service.Service) proto.MetricsCollectorGrpcServer {
	return &server{
		s: s,
	}
}

var emptyRes = &types.Empty{}

func (h *server) Send(stream proto.MetricsCollectorGrpc_SendServer) error {
	m := proto.MetricPool.Get().(*proto.Metric)
	defer func() {
		m.Values.Reset()
		m.Reset()
		proto.MetricPool.Put(m)
	}()
	for {
		if err := stream.RecvMsg(m); err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
			break
		}
		h.s.Send(m)
	}

	return nil
}
