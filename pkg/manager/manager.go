/*
Copyright 2019 The Kubernetes Authors.

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

package manager

import (
	goctx "context"
	"fmt"
	"os"

	"github.com/pkg/errors"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha4"
	bootstrapv1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1alpha4"
	ctrl "sigs.k8s.io/controller-runtime"

	infrav1a3 "sigs.k8s.io/cluster-api-provider-vsphere/api/v1alpha3"
	infrav1a4 "sigs.k8s.io/cluster-api-provider-vsphere/api/v1alpha4"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/record"
)

// Manager is a CAPV controller manager.
type Manager interface {
	ctrl.Manager

	// GetContext returns the controller manager's context.
	GetContext() *context.ControllerManagerContext
}

// New returns a new CAPV controller manager.
func New(opts Options) (Manager, error) {

	// Ensure the default options are set.
	opts.defaults()

	_ = clientgoscheme.AddToScheme(opts.Scheme)
	_ = clusterv1.AddToScheme(opts.Scheme)
	_ = infrav1a3.AddToScheme(opts.Scheme)
	_ = infrav1a4.AddToScheme(opts.Scheme)
	_ = bootstrapv1.AddToScheme(opts.Scheme)
	// +kubebuilder:scaffold:scheme

	podName, err := os.Hostname()
	if err != nil {
		podName = DefaultPodName
	}

	// Build the controller manager.
	mgr, err := ctrl.NewManager(opts.KubeConfig, ctrl.Options{
		Scheme:                  opts.Scheme,
		MetricsBindAddress:      opts.MetricsAddr,
		LeaderElection:          opts.LeaderElectionEnabled,
		LeaderElectionID:        opts.LeaderElectionID,
		LeaderElectionNamespace: opts.LeaderElectionNamespace,
		SyncPeriod:              &opts.SyncPeriod,
		Namespace:               opts.WatchNamespace,
		NewCache:                opts.NewCache,
		Port:                    opts.WebhookPort,
		HealthProbeBindAddress:  opts.HealthAddr,
		CertDir:                 opts.CertDir,
	})
	if err != nil {
		return nil, errors.Wrap(err, "unable to create manager")
	}

	// Build the controller manager context.
	controllerManagerContext := &context.ControllerManagerContext{
		Context:                 goctx.Background(),
		Namespace:               opts.WatchNamespace,
		Name:                    opts.PodName,
		LeaderElectionID:        opts.LeaderElectionID,
		LeaderElectionNamespace: opts.LeaderElectionNamespace,
		MaxConcurrentReconciles: opts.MaxConcurrentReconciles,
		Client:                  mgr.GetClient(),
		Logger:                  opts.Logger.WithName(opts.PodName),
		Recorder:                record.New(mgr.GetEventRecorderFor(fmt.Sprintf("%s/%s", opts.PodNamespace, podName))),
		Scheme:                  opts.Scheme,
		Username:                opts.Username,
		Password:                opts.Password,
	}

	// Add the requested items to the manager.
	if err := opts.AddToManager(controllerManagerContext, mgr); err != nil {
		return nil, errors.Wrap(err, "failed to add resources to the manager")
	}

	// +kubebuilder:scaffold:builder

	return &manager{
		Manager: mgr,
		ctx:     controllerManagerContext,
	}, nil
}

type manager struct {
	ctrl.Manager
	ctx *context.ControllerManagerContext
}

func (m *manager) GetContext() *context.ControllerManagerContext {
	return m.ctx
}
