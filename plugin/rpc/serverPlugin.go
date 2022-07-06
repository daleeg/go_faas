package main

import (
	"context"
	"errors"
	"flag"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
	"go_faas/plugin/rpc/pb"
	"google.golang.org/grpc"
	"log"
)

var (
	PackageName = "rpc"
)

type Config struct {
	zrpc.RpcServerConf
}

func PluginRpcServer(cfgFileName string) error {
	flag.Parse()

	var cfg Config
	conf.MustLoad(cfgFileName, &cfg)

	srv, err := zrpc.NewServer(cfg.RpcServerConf, func(s *grpc.Server) {
		pb.RegisterGreeterServer(s, &Hello{})
	})
	if err != nil {
		log.Fatal(err)
		return errors.New("start error")
	}
	srv.Start()
	return nil

}

type Hello struct{}

func (h *Hello) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "hello " + in.Name}, nil
}
