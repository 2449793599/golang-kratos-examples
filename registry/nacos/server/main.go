package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/go-kratos/examples/helloworld/helloworld"

	"github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	// NACOS相关依赖库
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: fmt.Sprintf("Welcome %+v!", in.Name)}, nil
}

func main() {

	// *****************************************************************************************
	sc := []constant.ServerConfig{ // NACOS服务器信息
		*constant.NewServerConfig("127.0.0.1", 8848),
	}

	cc := constant.ClientConfig{ // 当前服务注册信息
		NamespaceId:         "public",           // PUBLIC命名空间可以直接填写空字符串
		TimeoutMs:           5000,               // 请求NACOS服务超时时间
		NotLoadCacheAtStart: true,               // 不在启动时加载缓存
		LogDir:              "/tmp/nacos/log",   // 本地日志（当前盘符根目录下）
		CacheDir:            "/tmp/nacos/cache", // 本地缓存（当前盘符根目录下）
		RotateTime:          "1h",               // 本地日志的滚动时长
		MaxAge:              3,                  // 本地日志的最大保存时间
		LogLevel:            "debug",            // 本地日志级别
	}

	client, err := clients.NewNamingClient( // 创建NACOS客户端
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	if err != nil {
		log.Panic(err)
	}

	// *****************************************************************************************
	httpSrv := http.NewServer(
		http.Address(":8000"),
		http.Middleware(
			recovery.Recovery(),
		),
	)
	grpcSrv := grpc.NewServer(
		grpc.Address(":9000"),
		grpc.Middleware(
			recovery.Recovery(),
		),
	)

	s := &server{}

	pb.RegisterGreeterServer(grpcSrv, s)
	pb.RegisterGreeterHTTPServer(httpSrv, s)

	r := nacos.New(client) // KRATOS提供的基于NACOS的服务注册

	app := kratos.New(
		kratos.Name("helloworld"),
		kratos.Server(
			httpSrv,
			grpcSrv,
		),
		kratos.Registrar(r), // 将NACOS注册组件添加到KRATOS应用中
	)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}

}
