package main

import (
	"context"
	"fmt"
	"github.com/gogo/protobuf/types"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/maxim-kuderko/metrics-collector/internal/initializers"
	"github.com/maxim-kuderko/metrics-collector/internal/repositories"
	"github.com/maxim-kuderko/metrics-collector/internal/service"
	"github.com/maxim-kuderko/metrics-collector/proto"
	"github.com/spf13/viper"
	_ "github.com/un000/grpc-snappy"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
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
	serv := grpc.NewServer(grpc.MaxRecvMsgSize(1000 << 20))
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

func (h *server) Send(ctx context.Context, metrics *proto.MetricsRequest) (*empty.Empty, error) {
	for _, m := range metrics.Metric {
		h.s.Send(m)
	}
	return &empty.Empty{}, nil
}
