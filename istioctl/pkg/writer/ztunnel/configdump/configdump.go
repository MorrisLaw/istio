// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package configdump

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"

	"sigs.k8s.io/yaml"
)

// ConfigWriter is a writer for processing responses from the Ztunnel Admin config_dump endpoint
type ConfigWriter struct {
	Stdout      io.Writer
	ztunnelDump *ZtunnelDump
	FullDump    []byte
}

// Prime loads the config dump into the writer ready for printing
// func (c *ConfigWriter) Prime(b []byte) error {
// 	cd := ZtunnelDump{}
// 	// TODO(fisherxu): migrate this to jsonpb when issue fixed in golang
// 	// Issue to track -> https://github.com/golang/protobuf/issues/632
// 	err := json.Unmarshal(b, &cd)
// 	if err != nil {
// 		return fmt.Errorf("error unmarshalling config dump response from ztunnel: %v", err)
// 	}
// 	c.ztunnelDump = &cd
// 	return nil
// }

func (c *ConfigWriter) Prime(b []byte) error {
	fmt.Println("Raw JSON data:", string(b))
	cd := ZtunnelDump{}
	err := json.Unmarshal(b, &cd)
	if err != nil {
		return fmt.Errorf("error unmarshalling config dump response from ztunnel: %v", err)
	}
	fmt.Printf("Unmarshalled Data: %+v\n", cd)
	c.ztunnelDump = &cd
	return nil
}

// PrintBootstrapDump prints just the bootstrap config dump to the ConfigWriter stdout
func (c *ConfigWriter) PrintBootstrapDump(outputFormat string) error {
	// TODO
	return nil
}

func (c *ConfigWriter) PrintFullSummary() error {
	_, _ = c.Stdout.Write([]byte("\n"))
	if err := c.PrintWorkloadSummary(WorkloadFilter{Verbose: true}); err != nil {
		return err
	}
	_, _ = c.Stdout.Write([]byte("\n"))
	if err := c.PrintServiceSummary(ServiceFilter{}); err != nil {
		return err
	}
	_, _ = c.Stdout.Write([]byte("\n"))
	if err := c.PrintPolicySummary(PolicyFilter{}); err != nil {
		return err
	}
	_, _ = c.Stdout.Write([]byte("\n"))
	if err := c.PrintSecretSummary(); err != nil {
		return err
	}
	return nil
}

func (c *ConfigWriter) PrintFullDump(outputFormat string) error {
	out := c.FullDump
	if outputFormat == "yaml" {
		var err error
		out, err = yaml.JSONToYAML(out)
		if err != nil {
			return err
		}
	}
	fmt.Fprintln(c.Stdout, string(out))
	return nil
}

// PrintVersionSummary prints version information for Istio and Ztunnel from the config dump
func (c *ConfigWriter) PrintVersionSummary() error {
	// TODO
	return nil
}

// PrintPodRootCAFromDynamicSecretDump prints just pod's root ca from dynamic secret config dump to the ConfigWriter stdout
func (c *ConfigWriter) PrintPodRootCAFromDynamicSecretDump() (string, error) {
	// TODO
	return "", nil
}

func (c *ConfigWriter) tabwriter() *tabwriter.Writer {
	w := new(tabwriter.Writer).Init(c.Stdout, 0, 8, 1, ' ', 0)
	return w
}
