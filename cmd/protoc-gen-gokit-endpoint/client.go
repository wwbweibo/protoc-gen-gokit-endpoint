package main

import (
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
)

func generateClient(file *protogen.GeneratedFile, service *protogen.Service) {
	file.P(fmt.Sprintf("type %sGrpcClient struct {", service.GoName))
	file.P("\tconn *grpc.ClientConn")
	file.P("\tlogger log.Logger")
	file.P("\ttracer stdopentracing.Tracer")
	for _, method := range service.Methods {
		if method.Desc.IsStreamingServer() || method.Desc.IsStreamingClient() {
			continue
		}
		file.P(fmt.Sprintf("\t%sEndpoint endpoint.Endpoint", method.GoName))
	}
	file.P("}\n")
	generateNewClient(file, service)
}

func generateNewClient(file *protogen.GeneratedFile, service *protogen.Service) {
	file.P(fmt.Sprintf("func New%sGrpcClient(conn *grpc.ClientConn, logger log.Logger) *%sGrpcClient {", service.GoName, service.GoName))
	file.P(fmt.Sprintf("\treturn &%sGrpcClient{", service.GoName))
	file.P("\t\tconn: conn,")
	file.P("\t\tlogger: logger,")
	file.P("\t}")
	file.P("}")
	generateClientWithTracer(file, service)
	generateClientBuild(file, service)
	generateMethodImpl(file, service)
}

func generateClientBuild(file *protogen.GeneratedFile, service *protogen.Service) {
	file.P(fmt.Sprintf("func (client *%sGrpcClient) Build() {", service.GoName))
	for _, method := range service.Methods {
		if method.Desc.IsStreamingServer() || method.Desc.IsStreamingClient() {
			continue
		}
		file.P(fmt.Sprintf("\tclient.%sEndpoint = client.build%s()", method.GoName, method.GoName))
	}
	file.P("}")

	for _, method := range service.Methods {
		if method.Desc.IsStreamingServer() || method.Desc.IsStreamingClient() {
			continue
		}
		generateClientMethodBuild(file, service, method)
	}
}

func generateClientMethodBuild(file *protogen.GeneratedFile, service *protogen.Service, method *protogen.Method) {
	file.P(fmt.Sprintf("func (client *%sGrpcClient) build%s() endpoint.Endpoint {", service.GoName, method.GoName))
	file.P("\toptions := []grpctransport.ClientOption{}")
	file.P("\tif client.tracer != nil {")
	file.P("\t\toptions = append(options, grpctransport.ClientBefore(opentracing.ContextToGRPC(client.tracer, client.logger)))")
	file.P("\t}")
	file.P("\tendpoint := grpctransport.NewClient(")
	file.P("\t\tclient.conn,")
	file.P("\t\t\"" + service.Desc.FullName() + "\",")
	file.P("\t\t\"" + method.Desc.Name() + "\",")
	file.P("\t\tfunc(ctx context.Context,i interface{}) (interface{}, error) { return i, nil },")
	file.P("\t\tfunc(ctx context.Context,i interface{}) (interface{}, error) { return i, nil },")
	file.P("\t\t" + method.Output.Desc.Name() + "{}, options...).Endpoint()")
	file.P("\tif client.tracer != nil {")
	file.P("\t\tendpoint = opentracing.TraceClient(client.tracer, \"" + method.GoName + "\")(endpoint)")
	file.P("\t}")
	file.P("\treturn endpoint")
	file.P("}")
}

func generateMethodImpl(file *protogen.GeneratedFile, service *protogen.Service) {
	for _, method := range service.Methods {
		if method.Desc.IsStreamingServer() || method.Desc.IsStreamingClient() {
			continue
		}
		file.P(fmt.Sprintf("func (client *%sGrpcClient) %s(ctx context.Context, request *%s) (*%s, error) {",
			service.GoName, method.GoName, method.Input.GoIdent.GoName, method.Output.GoIdent.GoName))
		file.P(fmt.Sprintf("\tresp, err := client.%sEndpoint(ctx, request)", method.GoName))
		file.P("\tif err != nil {")
		file.P("\t\treturn nil, err")
		file.P("\t}")
		file.P(fmt.Sprintf("\treturn resp.(*%s), err", method.Output.GoIdent.GoName))
		file.P("}")
	}
}

func generateClientWithTracer(file *protogen.GeneratedFile, service *protogen.Service) {
	file.P(fmt.Sprintf("func (client *%sGrpcClient) WithTracing(tracer stdopentracing.Tracer) {", service.GoName))
	file.P("\tclient.tracer = tracer")
	file.P("}")
}
