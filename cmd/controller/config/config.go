/*
 * Copyright (c) 2018 VMware Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package config

import (
	"flag"
	"sync"

	"github.com/vmware/purser/pkg/client"
	group_client "github.com/vmware/purser/pkg/client/clientset/typed/groups/v1"
	subscriber_client "github.com/vmware/purser/pkg/client/clientset/typed/subscriber/v1"
	"github.com/vmware/purser/pkg/controller"
	"github.com/vmware/purser/pkg/controller/buffering"
	"github.com/vmware/purser/pkg/utils"
)

// Setup initialzes the controller configuration
func Setup(conf *controller.Config) {
	kubeconfig := flag.String("kubeconfig", "", "path to the kubeconfig file")
	flag.Parse()

	*conf = controller.Config{}
	conf.Kubeclient = utils.GetKubeclient(*kubeconfig)
	conf.Resource = controller.Resource{
		Pod:                   true,
		Node:                  true,
		PersistentVolume:      true,
		PersistentVolumeClaim: true,
		ReplicaSet:            true,
		Deployment:            true,
		StatefulSet:           true,
		DaemonSet:             true,
		Job:                   true,
		Service:               true,
	}
	conf.RingBuffer = &buffering.RingBuffer{Size: buffering.BufferSize, Mutex: &sync.Mutex{}}
	clientset, clusterConfig := client.GetAPIExtensionClient(*kubeconfig)
	conf.Groupcrdclient = group_client.NewGroupClient(clientset, clusterConfig)
	conf.Subscriberclient = subscriber_client.NewSubscriberClient(clientset, clusterConfig)
}
