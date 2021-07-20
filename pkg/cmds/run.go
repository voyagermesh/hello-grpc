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

package cmds

import (
	_ "net/http/pprof"

	"voyagermesh.dev/hello-grpc/pkg/cmds/server"
	_ "voyagermesh.dev/hello-grpc/pkg/hello"
	_ "voyagermesh.dev/hello-grpc/pkg/status"

	"github.com/spf13/cobra"
)

func NewCmdRun(stopCh <-chan struct{}) *cobra.Command {
	o := server.NewServerOptions()

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Launch Hello GRPC server",
		Long:  "Launch Hello GRPC server",
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Complete(); err != nil {
				return err
			}
			if err := o.Validate(args); err != nil {
				return err
			}
			if err := o.RunServer(stopCh); err != nil {
				return err
			}
			return nil
		},
	}

	flags := cmd.Flags()
	o.RecommendedOptions.AddFlags(flags)
	flags.BoolVar(&o.LogRPC, "log-rpc", o.LogRPC, "log RPC request and response")

	return cmd
}
