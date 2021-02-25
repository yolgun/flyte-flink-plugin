// Copyright 2021 Spotify AB.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package flink

import (
	"testing"

	flinkIdl "github.com/spotify/flyte-flink-plugin/gen/pb-go/flyteidl-flink"
	flyteIdlCore "github.com/lyft/flyteidl/gen/pb-go/flyteidl/core"
	"github.com/lyft/flyteplugins/go/tasks/pluginmachinery/core/mocks"
	"gotest.tools/assert"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestBuildFlinkClusterSpecValid(t *testing.T) {
	job := flinkIdl.FlinkJob{
		JarFile: "job.jar",
		FlinkProperties: map[string]string{
			"taskmanager.numberOfTaskSlots": "1",
		},
	}
	config := GetFlinkConfig()

	tID := &mocks.TaskExecutionID{}
	tID.OnGetID().Return(flyteIdlCore.TaskExecutionIdentifier{
		NodeExecutionId: &flyteIdlCore.NodeExecutionIdentifier{
			ExecutionId: &flyteIdlCore.WorkflowExecutionIdentifier{
				Name:    "name",
				Project: "project",
				Domain:  "domain",
			},
		},
	})
	tID.OnGetGeneratedName().Return("generated-name")

	taskCtx := &mocks.TaskExecutionMetadata{}
	taskCtx.OnGetTaskExecutionID().Return(tID)
	taskCtx.OnGetNamespace().Return("test-namespace")
	taskCtx.OnGetAnnotations().Return(make(map[string]string))
	taskCtx.OnGetLabels().Return(make(map[string]string))

	cluster, err := BuildFlinkClusterSpec(taskCtx, job, config)

	assert.NilError(t, err)
	assert.Equal(t, cluster.Spec.Image.Name, "flink-image")
}

func TestWithPersistentVolume(t *testing.T) {
	size := resource.MustParse("100Gi")
	job := flinkIdl.FlinkJob{
		JarFile: "job.jar",
		FlinkProperties: map[string]string{
			"taskmanager.numberOfTaskSlots": "1",
		},
		TaskManager: &flinkIdl.TaskManager{
			Resource: &flinkIdl.Resource{
				PersistentVolume: &flinkIdl.Resource_PersistentVolume{
					Type: flinkIdl.Resource_PersistentVolume_PD_SSD,
					Size: &size,
				},
			},
			Replicas: 1,
		},
	}
	config := GetFlinkConfig()

	tID := &mocks.TaskExecutionID{}
	tID.OnGetID().Return(flyteIdlCore.TaskExecutionIdentifier{
		NodeExecutionId: &flyteIdlCore.NodeExecutionIdentifier{
			ExecutionId: &flyteIdlCore.WorkflowExecutionIdentifier{
				Name:    "name",
				Project: "project",
				Domain:  "domain",
			},
		},
	})
	tID.OnGetGeneratedName().Return("generated-name")

	taskCtx := &mocks.TaskExecutionMetadata{}
	taskCtx.OnGetTaskExecutionID().Return(tID)
	taskCtx.OnGetNamespace().Return("test-namespace")
	taskCtx.OnGetAnnotations().Return(make(map[string]string))
	taskCtx.OnGetLabels().Return(make(map[string]string))

	cluster, err := BuildFlinkClusterSpec(taskCtx, job, config)

	assert.NilError(t, err)
	assert.Equal(t, cluster.Spec.Image.Name, "flink-image")
	assert.Assert(t, len(cluster.Spec.TaskManager.VolumeClaimTemplates) == 1)
	// it will include cache-colume and the required pv
	assert.Assert(t, len(cluster.Spec.TaskManager.VolumeMounts) == 2)
	assert.Assert(t, len(cluster.Spec.TaskManager.Volumes) == 2)
}

func TestBuildFlinkClusterSpecInvalid(t *testing.T) {
	job := flinkIdl.FlinkJob{
		FlinkProperties: map[string]string{
			"taskmanager.numberOfTaskSlots": "1",
		},
	}

	// Use empty config
	config := &Config{}

	tID := &mocks.TaskExecutionID{}
	tID.OnGetID().Return(flyteIdlCore.TaskExecutionIdentifier{
		NodeExecutionId: &flyteIdlCore.NodeExecutionIdentifier{
			ExecutionId: &flyteIdlCore.WorkflowExecutionIdentifier{
				Name:    "name",
				Project: "project",
				Domain:  "domain",
			},
		},
	})
	tID.OnGetGeneratedName().Return("generated-name")

	taskCtx := &mocks.TaskExecutionMetadata{}
	taskCtx.OnGetTaskExecutionID().Return(tID)
	taskCtx.OnGetNamespace().Return("test-namespace")
	taskCtx.OnGetAnnotations().Return(make(map[string]string))
	taskCtx.OnGetLabels().Return(make(map[string]string))

	_, err := BuildFlinkClusterSpec(taskCtx, job, config)

	assert.Error(t, err, "image name is unspecified")
}

func TestBuildFlinkClusterSpecServiceAccount(t *testing.T) {
	job := flinkIdl.FlinkJob{
		JarFile: "job.jar",
		FlinkProperties: map[string]string{
			"taskmanager.numberOfTaskSlots": "1",
		},
		ServiceAccount: "flink-user-service-account",
	}
	config := GetFlinkConfig()

	tID := &mocks.TaskExecutionID{}
	tID.OnGetID().Return(flyteIdlCore.TaskExecutionIdentifier{
		NodeExecutionId: &flyteIdlCore.NodeExecutionIdentifier{
			ExecutionId: &flyteIdlCore.WorkflowExecutionIdentifier{
				Name:    "name",
				Project: "project",
				Domain:  "domain",
			},
		},
	})
	tID.OnGetGeneratedName().Return("generated-name")

	taskCtx := &mocks.TaskExecutionMetadata{}
	taskCtx.OnGetTaskExecutionID().Return(tID)
	taskCtx.OnGetNamespace().Return("test-namespace")
	taskCtx.OnGetAnnotations().Return(make(map[string]string))
	taskCtx.OnGetLabels().Return(make(map[string]string))

	cluster, err := BuildFlinkClusterSpec(taskCtx, job, config)

	assert.NilError(t, err)
	assert.Equal(t, *cluster.Spec.ServiceAccountName, "flink-user-service-account")
}

func TestBuildFlinkClusterSpecImage(t *testing.T) {
	job := flinkIdl.FlinkJob{
		JarFile: "job.jar",
		FlinkProperties: map[string]string{
			"taskmanager.numberOfTaskSlots": "1",
		},
		Image: "flink-custom-image",
	}
	config := GetFlinkConfig()

	tID := &mocks.TaskExecutionID{}
	tID.OnGetID().Return(flyteIdlCore.TaskExecutionIdentifier{
		NodeExecutionId: &flyteIdlCore.NodeExecutionIdentifier{
			ExecutionId: &flyteIdlCore.WorkflowExecutionIdentifier{
				Name:    "name",
				Project: "project",
				Domain:  "domain",
			},
		},
	})
	tID.OnGetGeneratedName().Return("generated-name")

	taskCtx := &mocks.TaskExecutionMetadata{}
	taskCtx.OnGetTaskExecutionID().Return(tID)
	taskCtx.OnGetNamespace().Return("test-namespace")
	taskCtx.OnGetAnnotations().Return(make(map[string]string))
	taskCtx.OnGetLabels().Return(make(map[string]string))

	cluster, err := BuildFlinkClusterSpec(taskCtx, job, config)

	assert.NilError(t, err)
	assert.Equal(t, *&cluster.Spec.Image.Name, "flink-custom-image")
}
