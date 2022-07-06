package main

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
	"go_faas/plugin/rpc/pb"
	"log"
)

var (
	PackageName = "rpc"
)

func PluginRpcClient() error {
	client := zrpc.MustNewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: []string{"127.0.0.1:12379"},
			Key:   "hello.rpc",
		},
	})

	conn := client.Conn()
	hello := pb.NewGreeterClient(conn)
	reply, err := hello.SayHello(context.Background(), &pb.HelloRequest{Name: "go-zero"})
	if err != nil {
		log.Fatal(err)
		return errors.New("call error")
	}
	log.Println(reply.Message)
	return nil
}
