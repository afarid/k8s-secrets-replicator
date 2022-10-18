/*
Copyright 2022.

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
	replicasv1 "github.com/afarid/k8s-secrets-replicator/api/v1"
	v1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"regexp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"strings"
	"time"
)

// ReplicatedSecretReconciler reconciles a ReplicatedSecret object
type ReplicatedSecretReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=replicas.k8s.sh,resources=replicatedsecrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=replicas.k8s.sh,resources=replicatedsecrets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=replicas.k8s.sh,resources=replicatedsecrets/finalizers,verbs=update

//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=secrets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups="",resources=secrets/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ReplicatedSecret object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *ReplicatedSecretReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := log.FromContext(ctx).WithValues("ReplicatedSecret", req.Name)
	reqLogger.Info("Reconciling SnapshotSchedule")

	rs := &replicasv1.ReplicatedSecret{}
	err := r.Get(ctx, req.NamespacedName, rs)
	if err != nil {
		if kerrors.IsNotFound(err) {
			reqLogger.Info("object removed")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// eligibleNamespace is all namespaces in which the secret must be in
	var eligibleNamespaces = make(map[string]bool, 0)

	// Adding literal namespaces to list of eligibleNamespaces
	for _, namespace := range rs.Spec.Namespaces {
		eligibleNamespaces[namespace] = true
	}
	// Adding all existing namespaces with a specific prefix to eligibleNamespaces
	var allExisingNamespace v1.NamespaceList
	err = r.List(ctx, &allExisingNamespace)
	if err != nil {
		return ctrl.Result{}, err
	}
	for _, namespace := range allExisingNamespace.Items {
		if strings.HasPrefix(namespace.Name, rs.Spec.NamespacesPrefix) {
			eligibleNamespaces[namespace.Name] = true
		}
	}

	// Add all existing namespaces that match NamespacesRegex
	if len(rs.Spec.NamespacesRegex) > 0 {
		re, _ := regexp.Compile(rs.Spec.NamespacesRegex)
		for _, namespace := range allExisingNamespace.Items {
			if re.Match([]byte(namespace.Name)) {
				eligibleNamespaces[namespace.Name] = true
			}
		}
	}

	for namespace, _ := range eligibleNamespaces {
		secret := v1.Secret{}
		err = r.Get(ctx, types.NamespacedName{
			Namespace: namespace,
			Name:      req.Name,
		}, &secret)

		if err != nil {
			if kerrors.IsNotFound(err) {
				reqLogger.Info("Creating k8s secret in namespace", "namespace", namespace)
				secret.Data = rs.Spec.SecretData
				secret.ObjectMeta.Name = rs.Name
				secret.ObjectMeta.Namespace = namespace
				secret.Type = v1.SecretType(rs.Spec.Type)

				// Establish the parent-child relationship between my resource and the deployment
				if err = controllerutil.SetControllerReference(rs, &secret, r.Scheme); err != nil {
					reqLogger.Info("Failed to set replicated controller reference")
					return ctrl.Result{}, err
				}

				err = r.Create(ctx, &secret)
				if err != nil {
					return ctrl.Result{}, err
				}
				return ctrl.Result{}, nil
			}
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{
		RequeueAfter: time.Minute,
	}, nil
}

func (r ReplicatedSecretReconciler) createSecret() {

}

// SetupWithManager sets up the controller with the Manager.
func (r *ReplicatedSecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&replicasv1.ReplicatedSecret{}).Owns(&v1.Secret{}).
		Complete(r)
}
