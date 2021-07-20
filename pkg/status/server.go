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

package status

import (
	proto "voyagermesh.dev/hello-grpc/pkg/apis/status"
	"voyagermesh.dev/hello-grpc/pkg/cmds/server"

	"golang.org/x/net/context"
	v "gomodules.xyz/x/version"
)

func init() {
	server.GRPCEndpoints.Register(proto.RegisterStatusServiceServer, &Server{})
	server.GatewayEndpoints.Register(proto.RegisterStatusServiceHandlerFromEndpoint)
}

type Server struct {
	proto.UnimplementedStatusServiceServer
}

var _ proto.StatusServiceServer = &Server{}

func (s *Server) Status(ctx context.Context, req *proto.StatusRequest) (*proto.StatusResponse, error) {
	return &proto.StatusResponse{
		Version: &proto.Version{
			Version:         v.Version.Version,
			VersionStrategy: v.Version.VersionStrategy,
			CommitHash:      v.Version.CommitHash,
			GitBranch:       v.Version.GitBranch,
			GitTag:          v.Version.GitTag,
			CommitTimestamp: v.Version.CommitTimestamp,
		},
	}, nil
}
