/*
Copyright 2023.

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

package observability

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"sigs.k8s.io/yaml"

	// mimir "github.com/grafana/mimir/pkg/util/validation"

	"github.com/go-logr/logr"
	reconcilehelper "github.com/pluralsh/controller-reconcile-helper/pkg"
	observabilityv1alpha1 "github.com/pluralsh/trace-shield-controller/api/observability/v1alpha1"
	"github.com/pluralsh/trace-shield-controller/clients/keto"
)

// TenantReconciler reconciles a Tenant object
type TenantReconciler struct {
	client.Client
	KetoClient      *keto.KetoGrpcClient
	Scheme          *runtime.Scheme
	Config          *observabilityv1alpha1.Config
	mimirConfigData map[string]map[string]observabilityv1alpha1.MimirLimits
}

const (
	tenantFinalizerName = "tenants.observability.traceshield.io/finalizer"
)

//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch
//+kubebuilder:rbac:groups=observability.traceshield.io,resources=configs,verbs=get;list;watch
//+kubebuilder:rbac:groups=observability.traceshield.io,resources=tenants,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=observability.traceshield.io,resources=tenants/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=observability.traceshield.io,resources=tenants/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the 	1	`desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Tenant object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *TenantReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	tenantInstance := &observabilityv1alpha1.Tenant{}

	if err := r.Get(ctx, req.NamespacedName, tenantInstance); err != nil {
		if apierrs.IsNotFound(err) {
			// log.Info("Unable to fetch Tenant - skipping", "name", tenantInstance.Name)
			return ctrl.Result{}, nil
		}
		log.Error(err, "unable to fetch Tenant")
		return ctrl.Result{}, ignoreNotFound(err)
	}

	config := &observabilityv1alpha1.Config{}

	if err := r.Get(ctx, types.NamespacedName{Name: "config"}, config); err != nil {
		if apierrs.IsNotFound(err) {
			// log.Info("Unable to fetch Tenant - skipping", "name", tenantInstance.Name)
			return ctrl.Result{}, nil
		}
		log.Error(err, "unable to fetch Observability Config")
		return ctrl.Result{}, err
	}

	r.Config = config

	if err := r.getMimirConfigMap(ctx); err != nil {
		log.Error(err, "unable to fetch Mimir ConfigMap")
		return ctrl.Result{}, err
	}

	defer r.updateMimirConfigmap(ctx, log)

	// examine DeletionTimestamp to determine if object is under deletion
	if tenantInstance.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !controllerutil.ContainsFinalizer(tenantInstance, tenantFinalizerName) {
			controllerutil.AddFinalizer(tenantInstance, tenantFinalizerName)
			if err := r.Update(ctx, tenantInstance); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// The object is being deleted
		if controllerutil.ContainsFinalizer(tenantInstance, tenantFinalizerName) {
			// our finalizer is present, so lets handle any external dependency
			if err := r.deleteTenantResources(ctx, tenantInstance); err != nil {
				// if fail to delete the external dependency here, return with error
				// so that it can be retried
				log.Error(err, "unable to delete tenant resources", "name", tenantInstance.Name)
				return ctrl.Result{}, err
			}
			log.Info("deleted tenant resources", "name", tenantInstance.Name)

			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(tenantInstance, tenantFinalizerName)
			if err := r.Update(ctx, tenantInstance); err != nil {
				log.Error(err, "unable to remove finalizer", "name", tenantInstance.Name)
				return ctrl.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

	if err := r.KetoClient.CreateObservabilityTenantInKetoIfNotExists(ctx, tenantInstance.Name); err != nil {
		log.Error(err, "unable to create tenant in keto")
		return ctrl.Result{}, err
	}

	r.updateMimirConfigmapData(ctx, tenantInstance)

	return ctrl.Result{}, nil
}

func (r *TenantReconciler) updateMimirConfigmap(ctx context.Context, log logr.Logger) error {
	tenDat, _ := yaml.Marshal(r.mimirConfigData)

	configmapData := map[string]string{
		r.Config.Spec.Mimir.ConfigMap.Key: string(tenDat),
	}

	mimirConfigMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.Config.Spec.Mimir.ConfigMap.Name,
			Namespace: r.Config.Spec.Mimir.ConfigMap.Namespace,
		},
		Data: configmapData,
	}

	if err := reconcilehelper.ConfigMap(ctx, r.Client, mimirConfigMap, log); err != nil {
		log.Error(err, "Error reconciling ConfigMap", "name", mimirConfigMap.Name)
		return err
	}
	return nil
}

func (r *TenantReconciler) deleteTenantResources(ctx context.Context, tenant *observabilityv1alpha1.Tenant) error {
	delete(r.mimirConfigData["overrides"], tenant.Name)
	return r.KetoClient.DeleteObservabilityTenantInKeto(ctx, tenant.Name)
}

func (r *TenantReconciler) updateMimirConfigmapData(ctx context.Context, tenant *observabilityv1alpha1.Tenant) {
	if _, ok := r.mimirConfigData["overrides"]; !ok {
		r.mimirConfigData["overrides"] = make(map[string]observabilityv1alpha1.MimirLimits)
	}
	if tenant.Spec.Limits != nil && tenant.Spec.Limits.Mimir != nil {
		r.mimirConfigData["overrides"][tenant.Name] = *tenant.Spec.Limits.Mimir
	}
}

func (r *TenantReconciler) getMimirConfigMap(ctx context.Context) error {
	existingConfigmap := &corev1.ConfigMap{}

	currentTenantData := map[string]map[string]observabilityv1alpha1.MimirLimits{}

	err := r.Get(ctx, types.NamespacedName{Name: r.Config.Spec.Mimir.ConfigMap.Name, Namespace: r.Config.Spec.Mimir.ConfigMap.Namespace}, existingConfigmap)
	if err != nil {
		// TODO: handle error properly
		// if apierrs.IsNotFound(err) {
		// 	// log.Info("Unable to fetch Tenant - skipping", "name", tenantInstance.Name)
		// 	return ctrl.Result{}, nil
		// }
		// log.Error(err, "unable to fetch Tenant")
		// return ctrl.Result{}, ignoreNotFound(err)
	}

	if existingConfigmap.Data != nil {
		if tenantData, ok := existingConfigmap.Data[r.Config.Spec.Mimir.ConfigMap.Key]; ok {
			yaml.Unmarshal([]byte(tenantData), &currentTenantData)
		} else {
			// TODO: handle error properly
		}
	} else {
		// TODO: handle error properly
	}
	r.mimirConfigData = currentTenantData
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TenantReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&observabilityv1alpha1.Tenant{}).
		// WithEventFilter(predicate.Funcs{
		// 	CreateFunc: func(e event.CreateEvent) bool {

		// 		return true
		// 	},
		// 	UpdateFunc: func(e event.UpdateEvent) bool {
		// 		return true
		// 	},
		// 	DeleteFunc: func(e event.DeleteEvent) bool {
		// 		return true
		// 	},
		// }).
		Watches(
			&source.Kind{Type: &corev1.ConfigMap{}}, // TODO: change this watch to limit to the configmaps we care about
			handler.EnqueueRequestsFromMapFunc(r.findObjectsForConfigMap),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		Watches(
			&source.Kind{Type: &observabilityv1alpha1.Config{}},
			handler.EnqueueRequestsFromMapFunc(r.findObjectsForConfigMap),
		).
		Complete(r)
}

func (r *TenantReconciler) findObjectsForConfigMap(configMap client.Object) []reconcile.Request {
	if err := r.Get(context.TODO(), types.NamespacedName{Name: "config"}, r.Config); err != nil {
		if apierrs.IsNotFound(err) {
			// log.Info("Unable to fetch Tenant - skipping", "name", tenantInstance.Name)
			return []reconcile.Request{}
		}
		return []reconcile.Request{}
	}

	if configMap.GetName() == r.Config.Spec.Mimir.ConfigMap.Name && configMap.GetNamespace() == r.Config.Spec.Mimir.ConfigMap.Namespace { //TODO: expand to Loki and Mimir as well
		tenantList := &observabilityv1alpha1.TenantList{}
		err := r.List(context.TODO(), tenantList, &client.ListOptions{})
		if err != nil {
			return []reconcile.Request{}
		}

		requests := make([]reconcile.Request, len(tenantList.Items))
		for i, item := range tenantList.Items {
			requests[i] = reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      item.GetName(),
					Namespace: item.GetNamespace(),
				},
			}
		}
		return requests
	} else {
		return []reconcile.Request{}
	}
}

func ignoreNotFound(err error) error {
	if apierrs.IsNotFound(err) {
		return nil
	}
	return err
}
