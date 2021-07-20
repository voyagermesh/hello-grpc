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

package server

import (
	"net/http"
	"strings"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	gwrt "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	utilerrors "gomodules.xyz/errors"
	"gomodules.xyz/grpc-go-addons/endpoints"
	grpc_security "gomodules.xyz/grpc-go-addons/security"
	"gomodules.xyz/grpc-go-addons/server"
	"gomodules.xyz/grpc-go-addons/server/options"
	stringz "gomodules.xyz/x/strings"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	GRPCEndpoints    = endpoints.GRPCRegistry{}
	GatewayEndpoints = endpoints.ProxyRegistry{}
)

type ServerOptions struct {
	RecommendedOptions *options.RecommendedOptions
	LogRPC             bool
}

func NewServerOptions() *ServerOptions {
	o := &ServerOptions{
		RecommendedOptions: options.NewRecommendedOptions(),
	}
	return o
}

func (o ServerOptions) Validate(args []string) error {
	var errors []error
	errors = append(errors, o.RecommendedOptions.Validate()...)
	return utilerrors.NewAggregate(errors)
}

func (o *ServerOptions) Complete() error {
	return nil
}

func (o ServerOptions) Config() (*server.Config, error) {
	config := server.NewConfig()
	if err := o.RecommendedOptions.ApplyTo(config); err != nil {
		return nil, err
	}

	config.SetGRPCRegistry(GRPCEndpoints)
	config.SetProxyRegistry(GatewayEndpoints)

	optsGLog := []grpc_zap.Option{
		grpc_zap.WithDecider(func(methodFullName string, err error) bool {
			// will not log gRPC calls if it was a call to healthcheck and no error was raised
			if err == nil && methodFullName == "/voyagermesh.dev.hellogrpc.apis.status.StatusService/Status" {
				return false
			}

			// by default you will log all calls
			return o.LogRPC
		}),
	}
	payloadDecider := func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
		// will not log gRPC calls if it was a call to healthcheck and no error was raised
		if fullMethodName == "/voyagermesh.dev.hellogrpc.apis.status.StatusService/Status" {
			return false
		}

		// by default you will log all calls
		return o.LogRPC
	}

	zapLog, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	grpc_zap.ReplaceGrpcLogger(zapLog)

	config.GRPCServerOption(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_opentracing.StreamServerInterceptor(),
			grpc_prometheus.StreamServerInterceptor,
			grpc_zap.PayloadStreamServerInterceptor(zapLog, payloadDecider),
			grpc_zap.StreamServerInterceptor(zapLog, optsGLog...),
			grpc_security.StreamServerInterceptor(),
			grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_opentracing.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_zap.PayloadUnaryServerInterceptor(zapLog, payloadDecider),
			grpc_zap.UnaryServerInterceptor(zapLog, optsGLog...),
			grpc_security.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)

	config.GatewayMuxOption(gwrt.WithIncomingHeaderMatcher(func(h string) (string, bool) {
		if stringz.PrefixFold(h, "access-control-request-") ||
			stringz.PrefixFold(h, "k8s-") ||
			strings.EqualFold(h, "Origin") ||
			strings.EqualFold(h, "Cookie") ||
			strings.EqualFold(h, "X-Phabricator-Csrf") {
			return h, true
		}
		return "", false
	}),
		gwrt.WithOutgoingHeaderMatcher(func(h string) (string, bool) {
			if stringz.PrefixFold(h, "access-control-allow-") ||
				strings.EqualFold(h, "Set-Cookie") ||
				strings.EqualFold(h, "vary") ||
				strings.EqualFold(h, "x-content-type-options") ||
				stringz.PrefixFold(h, "x-ratelimit-") {
				return h, true
			}
			return "", false
		}),
		gwrt.WithMetadata(func(c context.Context, req *http.Request) metadata.MD {
			return metadata.Pairs("x-forwarded-method", req.Method)
		}),
		gwrt.WithErrorHandler(gwrt.DefaultHTTPErrorHandler))

	return config, nil
}

func (o ServerOptions) RunServer(stopCh <-chan struct{}) error {
	config, err := o.Config()
	if err != nil {
		return err
	}

	server, err := config.New()
	if err != nil {
		return err
	}

	return server.Run(stopCh)
}
