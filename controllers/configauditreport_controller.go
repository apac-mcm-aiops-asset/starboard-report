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
	"os/exec"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	aquasecurityv1alpha1 "github.com/aquasecurity/starboard/pkg/apis/aquasecurity/v1alpha1"
)

const (
	reportPath = "/report/"
)

// ConfigAuditReportReconciler reconciles a ConfigAuditReport object
type ConfigAuditReportReconciler struct {
	client.Client
	Scheme           *runtime.Scheme
	NamespaceWatched string
}

//+kubebuilder:rbac:groups=aquasecurity.github.io.my.domain,resources=configauditreports,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=aquasecurity.github.io.my.domain,resources=configauditreports/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=aquasecurity.github.io.my.domain,resources=configauditreports/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ConfigAuditReport object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *ConfigAuditReportReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// your logic here
	// Fetch the ConfigAuditReport instance
	instance := &aquasecurityv1alpha1.ConfigAuditReport{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	if !instance.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so do nothing
		return reconcile.Result{}, nil
	}

	// get owner of instance, maybe Replicaset, Pod .etc
	var ownerType, ownerName string
	for _, ref := range instance.GetOwnerReferences() {
		ownerType = ref.Kind
		ownerName = ref.Name
		break
	}

	command := r.buildCommand(ownerType, ownerName)
	logger.Info(command)
	// try to run the command here: "starboard get report deployment/nginx > nginx.deploy.html"
	cmd := exec.Command("sh", "-c", command)
	logger.Info("Exporting report and waiting for it to finish...")
	err = cmd.Run()
	if err != nil {
		logger.Error(err, "Exporting report finished with error")
		return reconcile.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ConfigAuditReportReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&aquasecurityv1alpha1.ConfigAuditReport{}).
		Complete(r)
}

// try to run the command here: "starboard get report deployment/nginx > nginx.deploy.html"
func (r *ConfigAuditReportReconciler) buildCommand(workloadType, workloadName string) string {
	return "starboard -n " + r.NamespaceWatched + " get report " + workloadType + "/" + workloadName + " > " + reportPath + workloadName + "." + workloadType + ".html"
}
