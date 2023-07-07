package main

import (
	"context"
	greeter "day01/proto"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	// pb "github.com/go-micro/examples/helloworld/proto"
	"github.com/go-micro/plugins/v4/registry/consul"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/server"
	"go-micro.dev/v4/transport"
	// proto "google.golang.org/protobuf/proto"
)

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, req *greeter.Request, rsp *greeter.Response) error {
	fmt.Println("he", req.GetName())
	rsp.Msg = "Hello1 " + req.Name
	return nil
}

type netListener struct{}

const (
	PORT int64 = 8081
	// rawListener WrapperListener = nil
)

func main() {

	reg := consul.NewRegistry(func(o *registry.Options) {
		o.Addrs = []string{"127.0.0.1:8500"} // consul服务注册服务发现地址
		o.Timeout = time.Second * 10
	})
	// 定制监听ipv4的地址
	rawListener, err := net.Listen("tcp4", ":"+strconv.FormatInt(PORT, 10))
	if err != nil {
		log.Panicln(err.Error())
	}
	ctx := context.WithValue(context.Background(), netListener{}, rawListener)

	// micro这个库有两层option，第一层是service本身的，service {opts[], once}
	// 第二层在service.opts.server.opts里面,server哪来的呢？都是newOptions函数自动生成的默认，
	// 然后再根据传入的参数改server的参数
	service := micro.NewService(
		// service.opts.Context!
		// micro.Context(ctx),
		// 每一个micro.xxxx返回都是一个符合Options接口的函数
		micro.Name("GreeterService"),
		micro.Version("v0.0.1"),
		micro.Address(":"+strconv.FormatInt(PORT, 10)),
		micro.Registry(reg),
		// service .opts.Server.opts.Context

		micro.AddListenOption(func(o *server.Options) {
			o.ListenOptions = append(o.ListenOptions,
				func(lo *transport.ListenOptions) { lo.Context = ctx },
			)
		}),
		micro.AddListenOption(func(o *server.Options) {
			o.Context = ctx
		}),
		// micro.AddListenOption(server.Advertise("127.0.0.0")),
		// micro.RegisterHandler()
	)

	// 这时候 s.opts 就有值了，但是里面的值都是函数，opts的值已经安装传入都改一遍了。

	// 配置中心
	// consulConfig, err := config.GetConsulConfig()
	// if err != nil {
	// 	logger.Fatal(err)
	// }

	// kongInfo, err := config.GetKongGatewayFromConsul(consulConfig, "kong-gateway-url-local")
	// if err != nil {
	// 	logger.Fatal(err)
	// }

	// fmt.Println("kong配置url：", kongInfo.Url)
	// // send to http://localhost:8001/services
	// var gw gateway.MyGateway = &gateway.KongGateway{Url: kongInfo.Url}
	// myip := gateway.GetMyIP()
	// myuri := myip + ":" + strconv.FormatInt(PORT, 10)
	// fmt.Println("注册IP及uri：", myuri)
	// err = gw.CreateService(myuri, "user") // 顺序调换一下
	// if err != nil {
	// 	log.Fatal("注册网关uri错误：", err)
	// }
	// err = gw.CreateRoutes("user")
	// if err != nil {
	// 	log.Fatal("注册网关路由错误：", err)
	// }

	// Mysql配置信息
	// mysqlInfo, err := config.GetMysqlFromConsul(consulConfig, "mysql")
	// if err != nil {
	// 	logger.Fatal(err)
	// }

	// logger.Info("Mysql配置信息:", mysqlInfo)

	service.Init()

	if err := greeter.RegisterGreeterHandler(service.Server(), new(Greeter)); err != nil {
		logger.Fatal(err)
	}

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
