[//]: # (desc: 演示使用goner/gin、goner/urllib、goner/otel/tracer/http + jaeger 做分布式链路追踪 )

<p>
    <a href="README.md">English</a>&nbsp ｜&nbsp 中文
</p>


##  项目搭建步骤
### 1. 创建服务段
```bash
mkdir server
cd server
go mod init examples/otel/tracer/server
```