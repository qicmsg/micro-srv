> 初始化项目
```
protoc -I. --micro_out=. --go_out=. proto/vessel/vessel.proto

go get -u github.com/micro/protobuf/proto
go get -u github.com/micro/protobuf/protoc-gen-go
# 安装微服务工具包
go get -u github.com/micro/micro

go mod download
go mod vendor
```

> 运行服务
```
make build && make run
```