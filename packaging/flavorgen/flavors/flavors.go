/*
Copyright 2020 The Kubernetes Authors.

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

package flavors

import (
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/cluster-api-provider-vsphere/packaging/flavorgen/flavors/crs"
	bootstrapv1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1alpha4"
)

func MultiNodeTemplateWithHAProxy() []runtime.Object {
	lb := newHAProxyLoadBalancer()
	vsphereCluster := newVSphereCluster(&lb)
	machineTemplate := newVSphereMachineTemplate()
	controlPlane := newKubeadmControlplane(444, machineTemplate, []bootstrapv1.File{})
	kubeadmJoinTemplate := newKubeadmConfigTemplate()
	cluster := newCluster(vsphereCluster, &controlPlane)
	machineDeployment := newMachineDeployment(cluster, machineTemplate, kubeadmJoinTemplate)
	clusterResourceSet := newClusterResourceSet(cluster)
	crsResourcesCSI := crs.CreateCrsResourceObjectsCSI(&clusterResourceSet)
	crsResourcesCPI := crs.CreateCrsResourceObjectsCPI(&clusterResourceSet)

	MultiNodeTemplate := []runtime.Object{
		&cluster,
		&lb,
		&vsphereCluster,
		&machineTemplate,
		&controlPlane,
		&kubeadmJoinTemplate,
		&machineDeployment,
		&clusterResourceSet,
	}

	MultiNodeTemplate = append(MultiNodeTemplate, crsResourcesCSI...)
	MultiNodeTemplate = append(MultiNodeTemplate, crsResourcesCPI...)

	return MultiNodeTemplate
}

func MultiNodeTemplateWithKubeVIP() []runtime.Object {
	vsphereCluster := newVSphereCluster(nil)
	machineTemplate := newVSphereMachineTemplate()
	controlPlane := newKubeadmControlplane(444, machineTemplate, newKubeVIPFiles())
	kubeadmJoinTemplate := newKubeadmConfigTemplate()
	cluster := newCluster(vsphereCluster, &controlPlane)
	machineDeployment := newMachineDeployment(cluster, machineTemplate, kubeadmJoinTemplate)
	clusterResourceSet := newClusterResourceSet(cluster)
	crsResourcesCSI := crs.CreateCrsResourceObjectsCSI(&clusterResourceSet)
	crsResourcesCPI := crs.CreateCrsResourceObjectsCPI(&clusterResourceSet)

	MultiNodeTemplate := []runtime.Object{
		&cluster,
		&vsphereCluster,
		&machineTemplate,
		&controlPlane,
		&kubeadmJoinTemplate,
		&machineDeployment,
		&clusterResourceSet,
	}

	MultiNodeTemplate = append(MultiNodeTemplate, crsResourcesCSI...)
	MultiNodeTemplate = append(MultiNodeTemplate, crsResourcesCPI...)

	return MultiNodeTemplate
}

func MultiNodeTemplateWithExternalLoadBalancer() []runtime.Object {
	vsphereCluster := newVSphereCluster(nil)
	machineTemplate := newVSphereMachineTemplate()
	controlPlane := newKubeadmControlplane(444, machineTemplate, nil)
	kubeadmJoinTemplate := newKubeadmConfigTemplate()
	cluster := newCluster(vsphereCluster, &controlPlane)
	machineDeployment := newMachineDeployment(cluster, machineTemplate, kubeadmJoinTemplate)
	clusterResourceSet := newClusterResourceSet(cluster)
	crsResourcesCSI := crs.CreateCrsResourceObjectsCSI(&clusterResourceSet)
	crsResourcesCPI := crs.CreateCrsResourceObjectsCPI(&clusterResourceSet)

	MultiNodeTemplate := []runtime.Object{
		&cluster,
		&vsphereCluster,
		&machineTemplate,
		&controlPlane,
		&kubeadmJoinTemplate,
		&machineDeployment,
		&clusterResourceSet,
	}
	MultiNodeTemplate = append(MultiNodeTemplate, crsResourcesCSI...)
	MultiNodeTemplate = append(MultiNodeTemplate, crsResourcesCPI...)

	return MultiNodeTemplate
}
