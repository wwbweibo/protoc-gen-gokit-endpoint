# protoc-gen-gokit-endpoint

protoc plugin used to generate go-kit grpc code 

安装

```bash
go install github.com/wwbweibo/protoc-gen-gokit-endpoint/cmd/protoc-gen-gokit-endpoint@latest
```

usage:
```bash
 protoc --proto_path=api/billing/v1 --go_out=. --go-grpc_out=. --gokit-endpoint_out=. package_service.proto quota_service.proto   query_service.proto
```

in your code 
```go
// here is your biz logic
packageService := usecase.PackageService{}
// this is the generated grpc server, you need pass your service into it
packageServer := v1.NewPackageService(packageService)

listener, err := net.Listen("tcp", ":8080")
if err != nil {
    panic(err)
}
grpcServer := grpc.NewServer()
// here to register your grpc server
v1.RegisterPackageServiceServer(grpcServer, packageServer)
// start serve
grpcServer.Serve(listener)
```