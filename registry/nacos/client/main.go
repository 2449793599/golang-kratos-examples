package main

import (
	"context"
	"log"

	"github.com/go-kratos/examples/helloworld/helloworld"
	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func main() {

	// *****************************************************************************************
	// 同服务端配置

	sc := []constant.ServerConfig{ // NACOS服务器信息
		*constant.NewServerConfig("127.0.0.1", 8848),
	}

	cc := &constant.ClientConfig{
		NamespaceId:         "public",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}

	cli, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  cc,
			ServerConfigs: sc,
		},
	)

	if err != nil {
		log.Panic(err)
	}

	// *****************************************************************************************
	conn, err := grpc.DialInsecure( // GRPC连接
		context.Background(),
		grpc.WithEndpoint("discovery:///helloworld.grpc"),
		grpc.WithDiscovery(nacos.New(cli)), // 通过服务发现执行
	)

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	client := helloworld.NewGreeterClient(conn)

	reply, err := client.SayHello(context.Background(), &helloworld.HelloRequest{Name: "kratos"})

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[grpc] SayHello %+v\n", reply)

}
