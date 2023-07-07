package main

import (
	"context"
	greeter "day01/proto"
	"fmt"

	"github.com/go-micro/plugins/v4/registry/consul"
	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/registry"
	// "golang.org/x/vuln/client"
	// proto "google.golang.org/protobuf/proto"
)

func main() {

	reg := consul.NewRegistry(func(o *registry.Options) {
		o.Addrs = []string{"127.0.0.1:8500"}
	})

	service := micro.NewService(
		// micro.Name("GreeterService"),
		// micro.Version("v0.0.1"),
		// micro.Address(":8080"),
		micro.Client(client.NewClient()),
		micro.Registry(reg),
		// micro.RegisterHandler()
	)

	service.Init()

	client := greeter.NewGreeterService("GreeterService", service.Client())
	fmt.Println("-")
	// 访问server
	rsp, err := client.Hello(context.Background(), &greeter.Request{Name: "fly Client"})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(rsp.Msg)
	}

}
