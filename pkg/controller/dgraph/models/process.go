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

package models

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/vmware/purser/pkg/controller/dgraph"
	"github.com/vmware/purser/pkg/controller/utils"
)

// Dgraph Model Constants
const (
	IsProc = "isProc"
)

// Proc schema in dgraph
type Proc struct {
	dgraph.ID
	IsProc    bool      `json:"isProc,omitemtpy"`
	Name      string    `json:"name,omitempty"`
	Interacts []*Pod    `json:"interacts,omitempty"`
	Container Container `json:"container,omitempty"`
	StartTime time.Time `json:"startTime,omitempty"`
	EndTime   time.Time `json:"endTime,omitempty"`
}

func newProc(procXID, procName, containerUID, containerXID string, creationTimeStamp time.Time) (*api.Assigned, error) {
	newProc := Proc{
		ID:        dgraph.ID{Xid: procXID},
		IsProc:    true,
		Name:      procName,
		Container: Container{ID: dgraph.ID{UID: containerUID, Xid: containerXID}},
		StartTime: creationTimeStamp,
	}
	bytes := utils.JSONMarshal(newProc)
	return dgraph.MutateNode(bytes)
}

// StoreProcess ...
func StoreProcess(procXID, procName string, podsXIDs []string, containerXID string,
	creationTimeStamp time.Time) error {
	containerUID := dgraph.GetUID(containerXID, IsContainer)
	if containerUID == "" {
		return fmt.Errorf("Container not persisted yet")
	}

	procUID := dgraph.GetUID(procXID, IsProc)
	if procUID == "" {
		assigned, err := newProc(procXID, procName, containerUID, containerXID, creationTimeStamp)
		if err != nil {
			logrus.Errorf("Unable to create proc: %s", procXID)
			return err
		}
		procUID = assigned.Uids["blank-0"]
	}

	pods := []*Pod{}
	for _, podXID := range podsXIDs {
		podUID := dgraph.GetUID(podXID, IsPod)
		if podUID != "" {
			pods = append(pods, &Pod{ID: dgraph.ID{UID: podUID, Xid: podXID}})
		}
	}

	updatedProc := Proc{
		ID:        dgraph.ID{UID: procUID, Xid: procXID},
		Interacts: pods,
	}
	bytes := utils.JSONMarshal(updatedProc)
	_, err := dgraph.MutateNode(bytes)
	return err
}

func deleteProcessesInTerminatedContainers(containers []*Container) error {
	procs := []Proc{}
	for _, container := range containers {
		for _, proc := range container.Procs {
			updatedProc := Proc{
				ID:      dgraph.ID{UID: proc.ID.UID, Xid: proc.ID.Xid},
				EndTime: container.EndTime,
			}
			procs = append(procs, updatedProc)
		}
	}
	bytes := utils.JSONMarshal(procs)
	_, err := dgraph.MutateNode(bytes)
	return err
}
