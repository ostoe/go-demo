package gateway

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
)

type KongGateway struct {
	Url string
}

type MyGateway interface {
	CreateService(svcURL, svcName string) error
	CreateRoutes(svcName string) error
}

// createService 在 Kong 上注册服务
func (k *KongGateway) CreateService(svcURL, svcName string) error {
	if !strings.HasPrefix(svcURL, "http://") { // 自动加上前缀'/'
		svcURL = "http://" + svcURL
	}
	data := url.Values{
		"name": {svcName},
		"url":  {svcURL},
	}
	fmt.Println("请求data：", data)

	resp, err := http.PostForm(k.Url, data)
	if err != nil {
		log.Println("http.PostForm() =>", err)
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("ioutil.ReadAll() =>", err)
		return err
	}

	fmt.Println("注册服务返回json：", string(body))
	return nil
}

// createRoutes 在 Kong 已注册的服务上注册路由
func (k *KongGateway) CreateRoutes(svcName string) error {
	// if !strings.HasPrefix(svcName, "/") { // 自动加上前缀'/'
	// 	svcName = svcName[1:]
	// }
	path := fmt.Sprintf("/%v", svcName)
	data := url.Values{
		"name":    {svcName},
		"paths[]": {path},
	}

	routesURL := fmt.Sprintf("%v/%v/routes", k.Url, svcName)
	fmt.Println("routesURL: ", routesURL)
	resp, err := http.PostForm(routesURL, data)
	if err != nil {
		log.Println("http.PostForm() =>", err)
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("ioutil.ReadAll() =>", err)
		return err
	}
	fmt.Println("注册路由返回", string(body))
	return nil
}

// mySvc 创建一个服务
func mySvc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello, You through Kong Gateway access my http service")
}

func GetMyIP() string {
	var myIP string
	adds, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, add := range adds {
		if ipnet, ok := add.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				myIP = ipnet.IP.String()
			}
		}
	}
	log.Printf("My IPv4 address is %v\n", myIP)
	return myIP
}
