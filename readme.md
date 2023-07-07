#TODO
 - config
 - 数据库
 - 中间件
 - 健康检查
 - 整合gin
 - 重启服务的更新机智，而不是create
 - 认证框架
 - http 和 grpc切换  http2，什么时候用http，网关入口流量是grpc还是http，在哪里做的转换

# 安装consul

启动 consul agent -data-dir=consul-tmp -dev  
consul agent -dev -data-dir=./deploy/consul-tmp
打开ui：
http://127.0.0.1:8500/ui/dc1/services


curl \
    --request PUT \
    --data 'hello consul' \
    http://127.0.0.1:8500/v1/kv/config

从yml放入
curl --request PUT --data-binary @config.yml http://localhost:8500/v1/kv/choice

config/

mysql

kong-gateway-url-local.json

{
  "url": "127.0.0.1"
  
}




# 安装kong

## docker启动

```bash
 docker network create kong-net
 ## 先去docker 上拉postgrs
  docker run -d --name kong-database \
  --network=kong-net \
  -p 5432:5432 \
  -e "POSTGRES_USER=kong" \
  -e "POSTGRES_DB=kong" \
  -e "POSTGRES_PASSWORD=kong123" \
  postgres:latest

docker run --rm --network=kong-net \
 -e "KONG_DATABASE=postgres" \
 -e "KONG_PG_HOST=kong-database" \
 -e "KONG_PG_PASSWORD=kong123" \
 -e "KONG_PASSWORD=kong123" \
kong/kong-gateway:3.3.0.0 kong migrations bootstrap

 docker run -d --name kong-gateway --network=kong-net  -e "KONG_DATABASE=postgres"  -e "KONG_PG_HOST=kong-database"  -e "KONG_PG_USER=kong"  -e "KONG_PG_PASSWORD=kongpass"  -e "KONG_PROXY_ACCESS_LOG=/dev/stdout"  -e "KONG_ADMIN_ACCESS_LOG=/dev/stdout"  -e "KONG_PROXY_ERROR_LOG=/dev/stderr"  -e "KONG_ADMIN_ERROR_LOG=/dev/stderr"  -e "KONG_ADMIN_LISTEN=0.0.0.0:8001"  -e "KONG_ADMIN_GUI_URL=http://localhost:8002"  -e KONG_LICENSE_DATA  -p 8000:8000  -p 8443:8443  -p 8001:8001  -p 8444:8444  -p 8002:8002  -p 8445:8445  -p 8003:8003  -p 8004:8004  kong/kong-gateway:3.3.0.0



```

## k8s启动
### kong gateway安装
### kong ingress gateway安装

## linux启动


# 两步走：

1. 创建Service：
curl -i -s -X POST http://localhost:8001/services \
  --data name=example_service \
  --data url='http://mockbin.org'

2. 创建路由
curl -i -X POST http://localhost:8001/services/example_service/routes \
  --data 'paths[]=/mock' \
  --data name=example_route

sidercar部署过程：

1. 启动consul server

2. 在每个服务进程节点，每个节点生成两个配置：

根据调用关系：
_services["hashicups-db.name"]="hashicups-db"
_services["hashicups-db.port"]="5432"
_services["hashicups-db.checks"]="hashicups-db:localhost:5432"
_services["hashicups-db.upstreams"]=""

_services["hashicups-api.name"]="hashicups-api"
_services["hashicups-api.port"]="8081"
_services["hashicups-api.checks"]="hashicups-api.public:localhost:8081,hashicups-api.product:localhost:9090,hashicups-api.payments:localhost:8080"
_services["hashicups-api.upstreams"]="hashicups-db:5432"

_services["hashicups-frontend.name"]="hashicups-frontend"
_services["hashicups-frontend.port"]="3000"
_services["hashicups-frontend.checks"]="hashicups-frontend:localhost:3000"
_services["hashicups-frontend.upstreams"]="hashicups-api:8081"

_services["hashicups-nginx.name"]="hashicups-nginx"
_services["hashicups-nginx.port"]="80"
_services["hashicups-nginx.checks"]="hashicups-nginx:localhost:80"
_services["hashicups-nginx.upstreams"]="hashicups-frontend:3000,hashicups-api:8081"

svc.服务网格.hcl

service {
  name = "sss1"
  id = "-1"
  tags = ["v1"]
  port = myservicePort
  connect {
    sidecar_service {
        proxy {
          upstreams = [
            一种是：
            {
              destination_name = "${_UPS_NAME}"
              local_bind_port = ${_UPS_PORT}
            }
            另一种是：

            ${_SERVICE_DEF_UPS}
          ]
        }
     }
  } 

  check 
}


svc.发现.hcl
service {
  name = "${_svc_name}"
  id = "${_svc_name}-1"
  tags = ["v1"]
  port = ${_svc_port}
  ${_svc_token}

  check
}

启动agent

consul members

Node                Address          Status  Type    Build   Protocol  DC   Partition  Segment
consul-server-0     10.0.4.189:8301  alive   server  1.15.2  2         dc1  default    <all>
hashicups-api       10.0.4.75:8301   alive   client  1.15.2  2         dc1  default    <default>
hashicups-db        10.0.4.217:8301  alive   client  1.15.2  2         dc1  default    <default>
hashicups-frontend  10.0.4.35:8301   alive   client  1.15.2  2         dc1  default    <default>
hashicups-nginx     10.0.4.59:8301   alive   client  1.15.2  2         dc1  default    <default>

consul进程 port：统一为 8301

![img](https://developer.hashicorp.com/_next/image?url=https%3A%2F%2Fcontent.hashicorp.com%2Fapi%2Fassets%3Fproduct%3Dtutorials%26version%3Dmain%26asset%3Dpublic%252Fimg%252Fconsul%252Fgetting-started-vms%252Fgs_vms-diagram-02.png%26width%3D2048%26height%3D622&w=3840&q=75)

3. 部署envoy sidercar
![img](https://developer.hashicorp.com/_next/image?url=https%3A%2F%2Fcontent.hashicorp.com%2Fapi%2Fassets%3Fproduct%3Dtutorials%26version%3Dmain%26asset%3Dpublic%252Fimg%252Fconsul%252Fgetting-started-vms%252Fgs_vms-diagram-03.png%26width%3D2048%26height%3D622&w=3840&q=75)

每个节点安装envoy，
生成intentions：缺省情况下，初始 Consul 配置会拒绝所有服务连接。我们 建议在生产环境中使用此设置以遵循“最低特权” 原则，除非明确定义，否则限制所有网络访问。有点acl的意思

4. reload consul，加载acl

5.  连接envoy：/usr/bin/consul connect envoy \
    -token=${CONSUL_HTTP_TOKEN} \
    -envoy-binary /usr/bin/envoy \
    -sidecar-for hashicups-db-1 > /tmp/sidecar-proxy.log 2>&1 &

6. 所有服务重启监听端口为：localhost

7. 目的：启用 Consul API 网关和 将其配置为允许通过端口 8443 访问
![img](https://developer.hashicorp.com/_next/image?url=https%3A%2F%2Fcontent.hashicorp.com%2Fapi%2Fassets%3Fproduct%3Dtutorials%26version%3Dmain%26asset%3Dpublic%252Fimg%252Fconsul%252Fgetting-started-vms%252Fgs_vms-diagram-03.png%26width%3D2048%26height%3D622&w=3840&q=75) ->

![img](https://developer.hashicorp.com/_next/image?url=https%3A%2F%2Fcontent.hashicorp.com%2Fapi%2Fassets%3Fproduct%3Dtutorials%26version%3Dmain%26asset%3Dpublic%252Fimg%252Fconsul%252Fgetting-started-vms%252Fgs_vms-diagram-04.png%26width%3D1280%26height%3D410&w=3840&q=75)

增加一个api-gateway节点：


8. 生成api网关规则：
optional：生成证书
9. 生成api网关路由

10。 跟之前一样：/usr/bin/consul connect envoy \
  -gateway api \
  -register \
  -service gateway-api \
  -token=${CONSUL_AGENT_TOKEN} \
  -envoy-binary /usr/bin/envoy > /tmp/api-gw-proxy.log 2>&1 &

11. 添加intention【acl allow】

12. 监控： Grafana
![img](https://developer.hashicorp.com/_next/image?url=https%3A%2F%2Fcontent.hashicorp.com%2Fapi%2Fassets%3Fproduct%3Dtutorials%26version%3Dmain%26asset%3Dpublic%252Fimg%252Fconsul%252Fgetting-started-vms%252Fgs_vms-diagram-05.png%26width%3D1178%26height%3D574&w=3840&q=75)

安装grafana-agent  ，配置 Grafana 代理，生成yml配置，启动grafana 从yml配置，




# k8s集群使用consul

1. 安装集群kind，
2. 使用helm安装consul，不仅仅是server进程，而且包含集群配置，自动注入sidecar到pod
启动服务……正常启动，额外添加一些consul标注：
```yml
consul.hashicorp.com/connect-inject: "true"，启用代理注入
consul.hashicorp.com/connect-service-upstreams: "frontend:3000, public-api:8080" 指定调用关系
```
3. 部署服务，正常部署就行，端口对应起来
4. 安装consul api-gateway，定制CRD，与k8s集成的，类似于nginx ingress 
5. 安装rbac，可以关联api调用权限
6. 安装监控服务，以envoy proxy的形式存在

![img](https://developer.hashicorp.com/_next/image?url=https%3A%2F%2Fcontent.hashicorp.com%2Fapi%2Fassets%3Fproduct%3Dtutorials%26version%3Dmain%26asset%3Dpublic%252Fimg%252Fconsul-splitting-architecture.png%26width%3D4561%26height%3D2914&w=3840&q=75)


7. 流量拆分，做config.json配置，根据标签配置权重
lua  http    nginx 进程    nginx.conf  
