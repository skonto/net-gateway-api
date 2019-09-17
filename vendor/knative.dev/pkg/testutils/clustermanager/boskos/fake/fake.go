/*
Copyright 2019 The Knative Authors

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

package fake

import (
	"fmt"

	boskoscommon "k8s.io/test-infra/boskos/common"
	"knative.dev/pkg/testutils/clustermanager/boskos"
)

// FakeBoskosClient implements boskos.Operation
type FakeBoskosClient struct {
	resources []*boskoscommon.Resource
}

// AcquireGKEProject fakes to be no op
func (c *FakeBoskosClient) AcquireGKEProject(host *string) (*boskoscommon.Resource, error) {
	for _, res := range c.resources {
		if res.State == boskoscommon.Free {
			res.State = boskoscommon.Busy
			return res, nil
		}
	}
	return nil, fmt.Errorf("no GKE project available")
}

// ReleaseGKEProject fakes to be no op
func (c *FakeBoskosClient) ReleaseGKEProject(host *string, name string) error {
	if nil == host {
		return fmt.Errorf("host has to be set")
	}

	for _, res := range c.resources {
		if res.Name == name {
			if res.Owner == *host {
				res.Owner = ""
				res.State = boskoscommon.Free
				return nil
			} else {
				return fmt.Errorf("Got owner: '%s', expect owner: '%s'", res.Owner, *host)
			}
		}
	}
	return fmt.Errorf("resource doesn't exist yet: '%s'", name)
}

// NewGKEProject adds Boskos resources for testing purpose
func (c *FakeBoskosClient) NewGKEProject(name string) {
	c.resources = append(c.resources, &boskoscommon.Resource{
		Type:  boskos.GKEProjectResource,
		Name:  name,
		State: boskoscommon.Free,
	})
}