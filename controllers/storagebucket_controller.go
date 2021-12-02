/*
Copyright 2021.

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

package controllers

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	demov1 "github.com/hhiroshell/storage-bucket-prober/api/v1"
)

// StorageBucketReconciler reconciles a StorageBucket object
type StorageBucketReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	ProbeInterval time.Duration
}

//+kubebuilder:rbac:groups=demo.hhiroshell.github.com,resources=storagebuckets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=demo.hhiroshell.github.com,resources=storagebuckets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=demo.hhiroshell.github.com,resources=storagebuckets/finalizers,verbs=update

func (r *StorageBucketReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("probe")

	var storageBucket demov1.StorageBucket
	if err := r.Client.Get(ctx, req.NamespacedName, &storageBucket); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	now := metav1.Now()
	storageBucket.Status.LastProbeTime = &now

	// dummy
	storageBucket.Status.Available = true

	if err := r.Client.Status().Patch(ctx, &storageBucket, client.Merge); err != nil {
		return ctrl.Result{Requeue: false}, err
	}

	return ctrl.Result{}, nil
}

type ticker struct {
	events chan event.GenericEvent

	interval time.Duration
}

func (t *ticker) Start(ctx context.Context) error {
	ticker := time.NewTicker(t.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			t.events <- event.GenericEvent{}
		}
	}
}

func (t *ticker) NeedLeaderElection() bool {
	return true
}

func (r *StorageBucketReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager) error {
	// Register the ticker as a runnable managed by the manager.
	// The storage bucket controller's reconciliation will be triggered by the events send by the ticker.
	events := make(chan event.GenericEvent)
	err := mgr.Add(&ticker{
		events:   events,
		interval: r.ProbeInterval,
	})
	if err != nil {
		return err
	}

	// The storage bucket controller watches the events send by the ticker.
	// When it receives an event from ticker, reconcile requests for all storage buckets will be enqueued (= polling).
	source := source.Channel{
		Source:         events,
		DestBufferSize: 0,
	}

	handler := handler.EnqueueRequestsFromMapFunc(func(object client.Object) []reconcile.Request {
		storageBuckets := demov1.StorageBucketList{}
		mgr.GetCache().List(ctx, &storageBuckets)

		var requests []reconcile.Request
		for _, storageBucket := range storageBuckets.Items {
			requests = append(requests, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      storageBucket.Name,
					Namespace: storageBucket.Namespace,
				},
			})
		}

		return requests
	})

	return ctrl.NewControllerManagedBy(mgr).
		For(&demov1.StorageBucket{}).
		Watches(&source, handler).
		Complete(r)
}
