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
	"strings"

	cache "github.com/patrickmn/go-cache"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	aquasecurityv1alpha1 "github.com/aquasecurity/starboard/pkg/apis/aquasecurity/v1alpha1"
)

const (
	reportPath    = "/report/"
	finalizerName = "aquasecurity.starboard/finalizer"
)

// ConfigAuditReportReconciler reconciles a ConfigAuditReport object
type ConfigAuditReportReconciler struct {
	client.Client
	Scheme           *runtime.Scheme
	NamespaceWatched string
	Cache            *cache.Cache
}

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

	workloadInfo, err := r.getWorkloadInfo(instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	// if !instance.ObjectMeta.DeletionTimestamp.IsZero() {
	// 	// The object is not being deleted, so do nothing
	// 	logger.Info("Remove the original report")
	// 	logger.Info("workloadInfo")
	// 	err := r.removeReport(workloadInfo)
	// 	if err != nil {
	// 		return reconcile.Result{}, err
	// 	}
	// 	return reconcile.Result{}, nil
	// }

	// examine DeletionTimestamp to determine if object is under deletion
	if instance.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !containsString(instance.GetFinalizers(), finalizerName) {
			addFinalizer(instance, finalizerName)
			if err := r.Client.Update(context.TODO(), instance); err != nil {
				return reconcile.Result{}, err
			}
		}
	} else {
		// The object is being deleted
		if containsString(instance.GetFinalizers(), finalizerName) {
			// our finalizer is present, so lets handle any external dependency
			if err := r.removeReport(workloadInfo); err != nil {
				// if fail to delete the external dependency here, return with error
				// so that it can be retried
				return reconcile.Result{}, err
			}

			// remove our finalizer from the list and update it.
			removeFinalizer(instance, finalizerName)
			if err := r.Client.Update(context.TODO(), instance); err != nil {
				return reconcile.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return reconcile.Result{}, nil
	}

	err = r.generateReport(ctx, workloadInfo)
	if err != nil {
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

func (r *ConfigAuditReportReconciler) getWorkloadInfo(instance *aquasecurityv1alpha1.ConfigAuditReport) (string, error) {
	itemKey := instance.Namespace + "/" + instance.Name
	workloadInfo, found := r.Cache.Get(itemKey)
	if found {
		return workloadInfo.(string), nil
	}

	ownerType, ownerName, err := r.findOwner(instance)
	if err != nil {
		return "", err
	}

	workloadInfo = ownerType + "|" + ownerName
	r.Cache.Set(itemKey, workloadInfo, cache.NoExpiration)

	return workloadInfo.(string), nil
}

func (r *ConfigAuditReportReconciler) generateReport(ctx context.Context, workloadInfo string) error {
	logger := log.FromContext(ctx)

	command := r.buildCommand(workloadInfo)
	logger.Info(command)
	// try to run the command here: "starboard get report deployment/nginx > nginx.deploy.html"
	cmd := exec.Command("sh", "-c", command)
	logger.Info("Exporting report and waiting for it to finish...")
	err := cmd.Run()
	if err != nil {
		logger.Error(err, "Exporting report finished with error")
		return err
	}

	return nil
}

func (r *ConfigAuditReportReconciler) removeReport(workloadInfo string) error {
	// workloadInfos := strings.Split(workloadInfo, "|")
	// workloadType := workloadInfos[0]
	// workloadName := workloadInfos[1]

	// err := os.Remove(buildReportName(workloadType, workloadName))
	// if err != nil {
	// 	return err
	// }

	return nil
}

// try to run the command here: "starboard get report deployment/nginx > nginx.deploy.html"
func (r *ConfigAuditReportReconciler) buildCommand(workloadInfo string) string {
	workloadInfos := strings.Split(workloadInfo, "|")
	workloadType := workloadInfos[0]
	workloadName := workloadInfos[1]

	return "starboard -n " + r.NamespaceWatched + " get report " + workloadType + "/" + workloadName + " > " + buildReportName(workloadType, workloadName)
}

func buildReportName(workloadType, workloadName string) string {
	return reportPath + workloadName + "." + workloadType + ".html"
}

func (r *ConfigAuditReportReconciler) findOwner(instance *aquasecurityv1alpha1.ConfigAuditReport) (string, string, error) {
	// get owner of instance, maybe Replicaset, Pod .etc
	var ownerType, ownerName string
	for _, ref := range instance.GetOwnerReferences() {
		ownerType = ref.Kind
		ownerName = ref.Name
		break
	}

	// if report's owner is ReplicaSet, then found the owner of the ReplicaSet, the Deployment
	if ownerType == "ReplicaSet" {
		nameNamespace := types.NamespacedName{
			Name:      ownerName,
			Namespace: instance.Namespace,
		}

		ownerReplicaSet := &v1.ReplicaSet{}
		err := r.Client.Get(context.TODO(), nameNamespace, ownerReplicaSet)
		if err != nil {
			return "", "", err
		}

		for _, ref := range ownerReplicaSet.GetOwnerReferences() {
			ownerType = ref.Kind
			ownerName = ref.Name
			break
		}
	}

	return ownerType, ownerName, nil
}

// Helper functions to check and remove string from a slice of strings.
func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}

// AddFinalizer accepts an Object and adds the provided finalizer if not present.
func addFinalizer(o *aquasecurityv1alpha1.ConfigAuditReport, finalizer string) {
	f := o.GetFinalizers()
	for _, e := range f {
		if e == finalizer {
			return
		}
	}
	o.SetFinalizers(append(f, finalizer))
}

// RemoveFinalizer accepts an Object and removes the provided finalizer if present.
func removeFinalizer(o *aquasecurityv1alpha1.ConfigAuditReport, finalizer string) {
	f := o.GetFinalizers()
	for i := 0; i < len(f); i++ {
		if f[i] == finalizer {
			f = append(f[:i], f[i+1:]...)
			i--
		}
	}
	o.SetFinalizers(f)
}
