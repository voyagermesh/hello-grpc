/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package hello

import (
	"fmt"
	"time"

	proto "voyagermesh.dev/hello-grpc/pkg/apis/hello/v1alpha1"
	"voyagermesh.dev/hello-grpc/pkg/cmds/server"

	"golang.org/x/net/context"
)

func init() {
	server.GRPCEndpoints.Register(proto.RegisterHelloServiceServer, &Server{})
	server.GatewayEndpoints.Register(proto.RegisterHelloServiceHandlerFromEndpoint)
}

type Server struct {
	proto.UnimplementedHelloServiceServer
}

var _ proto.HelloServiceServer = &Server{}

func (s *Server) Intro(ctx context.Context, req *proto.IntroRequest) (*proto.IntroResponse, error) {
	return &proto.IntroResponse{
		Intro: fmt.Sprintf("hello, %s!", req.Name),
	}, nil
}

func (s *Server) Stream(req *proto.IntroRequest, stream proto.HelloService_StreamServer) error {
	for i := 0; i < 60; i++ {
		intro := fmt.Sprintf("%d: hello, %s!", i, req.Name)
		if err := stream.Send(&proto.IntroResponse{Intro: intro}); err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}
