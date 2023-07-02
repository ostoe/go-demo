package main

import (
	"context"
	"day01/config"
	gateway "day01/mygateway"
	greeter "day01/proto"
	"fmt"
	"log"
	"strconv"

	// pb "github.com/go-micro/examples/helloworld/proto"
	"github.com/go-micro/plugins/v4/registry/consul"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/registry"
	// proto "google.golang.org/protobuf/proto"
)

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, req *greeter.Request, rsp *greeter.Response) error {
	fmt.Println("he", req.GetName())
	rsp.Msg = "Hello " + req.Name
	return nil
}

const (
	PORT int64 = 8080
)

func main() {

	reg := consul.NewRegistry(func(o *registry.Options) {
		o.Addrs = []string{"127.0.0.1:8500"} // consul服务注册服务发现地址
	})

	service := micro.NewService(
		micro.Name("GreeterService"),
		micro.Version("v0.0.1"),
		micro.Address(":"+strconv.FormatInt(PORT, 10)),
		micro.Registry(reg),
		// micro.RegisterHandler()
	)

	// 配置中心
	consulConfig, err := config.GetConsulConfig()
	if err != nil {
		logger.Fatal(err)
	}

	kongInfo, err := config.GetKongGatewayFromConsul(consulConfig, "kong-gateway-url-local")
	if err != nil {
		logger.Fatal(err)
	}

	fmt.Println("kong配置url：", kongInfo.Url)
	// send to http://localhost:8001/services
	var gw gateway.MyGateway = &gateway.KongGateway{Url: kongInfo.Url}
	myip := gateway.GetMyIP()
	myuri := myip + ":" + strconv.FormatInt(PORT, 10)
	fmt.Println("注册IP及uri：", myuri)
	err = gw.CreateService(myuri, "user") // 顺序调换一下
	if err != nil {
		log.Fatal("注册网关uri错误：", err)
	}
	err = gw.CreateRoutes("user")
	if err != nil {
		log.Fatal("注册网关路由错误：", err)
	}

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
