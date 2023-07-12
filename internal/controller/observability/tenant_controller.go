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

	"sigs.k8s.io/yaml"

	// mimir "github.com/grafana/mimir/pkg/util/validation"

	"github.com/go-logr/logr"
	reconcilehelper "github.com/pluralsh/controller-reconcile-helper/pkg/reconcile-helper/core"
	observabilityv1alpha1 "github.com/pluralsh/trace-shield-controller/api/observability/v1alpha1"
	"github.com/pluralsh/trace-shield-controller/clients/keto"
)

// TenantReconciler reconciles a Tenant object
type TenantReconciler struct {
	client.Client
	KetoClient      *keto.KetoGrpcClient
	Scheme          *runtime.Scheme
	Config          *observabilityv1alpha1.Config
	mimirConfigData mimirConfigData
	lokiConfigData  lokiConfigData
	tempoConfigData tempoConfigData
}

type mimirConfigData struct {
	Overrides                             map[string]observabilityv1alpha1.MimirLimits `yaml:"overrides" json:"overrides"`
	observabilityv1alpha1.MimirConfigSpec `yaml:",inline"`
}

type lokiConfigData struct {
	Overrides                            map[string]observabilityv1alpha1.LokiLimits `yaml:"overrides" json:"overrides"`
	observabilityv1alpha1.LokiConfigSpec `yaml:",inline"`
}

type tempoConfigData struct {
	Overrides map[string]observabilityv1alpha1.TempoLimits `yaml:"overrides" json:"overrides"`
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

	// TODO: these should be optional and skipped if not configured
	if err := r.getMimirConfigMap(ctx); err != nil {
		log.Error(err, "unable to fetch Mimir ConfigMap")
		return ctrl.Result{}, err
	}

	if err := r.getLokiConfigMap(ctx); err != nil {
		log.Error(err, "unable to fetch Loki ConfigMap")
		return ctrl.Result{}, err
	}

	if err := r.getTempoConfigMap(ctx); err != nil {
		log.Error(err, "unable to fetch Tempo ConfigMap")
		return ctrl.Result{}, err
	}

	defer r.updateMimirConfigmap(ctx, log)
	defer r.updateLokiConfigmap(ctx, log)
	defer r.updateTempoConfigmap(ctx, log)

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
	r.updateLokiConfigmapData(ctx, tenantInstance)
	r.updateTempoConfigmapData(ctx, tenantInstance)

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

func (r *TenantReconciler) updateLokiConfigmap(ctx context.Context, log logr.Logger) error {
	tenDat, _ := yaml.Marshal(r.lokiConfigData)

	configmapData := map[string]string{
		r.Config.Spec.Loki.ConfigMap.Key: string(tenDat),
	}

	lokiConfigMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.Config.Spec.Loki.ConfigMap.Name,
			Namespace: r.Config.Spec.Loki.ConfigMap.Namespace,
		},
		Data: configmapData,
	}

	if err := reconcilehelper.ConfigMap(ctx, r.Client, lokiConfigMap, log); err != nil {
		log.Error(err, "Error reconciling ConfigMap", "name", lokiConfigMap.Name)
		return err
	}
	return nil
}

func (r *TenantReconciler) updateTempoConfigmap(ctx context.Context, log logr.Logger) error {
	tenDat, _ := yaml.Marshal(r.tempoConfigData)

	configmapData := map[string]string{
		r.Config.Spec.Tempo.ConfigMap.Key: string(tenDat),
	}

	tempoConfigMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.Config.Spec.Tempo.ConfigMap.Name,
			Namespace: r.Config.Spec.Tempo.ConfigMap.Namespace,
		},
		Data: configmapData,
	}

	if err := reconcilehelper.ConfigMap(ctx, r.Client, tempoConfigMap, log); err != nil {
		log.Error(err, "Error reconciling ConfigMap", "name", tempoConfigMap.Name)
		return err
	}
	return nil
}

func (r *TenantReconciler) deleteTenantResources(ctx context.Context, tenant *observabilityv1alpha1.Tenant) error {
	delete(r.mimirConfigData.Overrides, tenant.Name)
	delete(r.lokiConfigData.Overrides, tenant.Name)
	delete(r.tempoConfigData.Overrides, tenant.Name)
	return r.KetoClient.DeleteObservabilityTenantInKeto(ctx, tenant.Name)
}

func (r *TenantReconciler) updateMimirConfigmapData(ctx context.Context, tenant *observabilityv1alpha1.Tenant) {
	if tenant.Spec.Limits != nil && tenant.Spec.Limits.Mimir != nil {
		r.mimirConfigData.Overrides[tenant.Name] = *tenant.Spec.Limits.Mimir
	}
	// update the global mimir config
	if r.Config.Spec.Mimir.Config != nil {
		r.mimirConfigData.MimirConfigSpec = *r.Config.Spec.Mimir.Config
	} else {
		r.mimirConfigData.MimirConfigSpec = observabilityv1alpha1.MimirConfigSpec{}
	}
}

func (r *TenantReconciler) updateLokiConfigmapData(ctx context.Context, tenant *observabilityv1alpha1.Tenant) {
	if tenant.Spec.Limits != nil && tenant.Spec.Limits.Loki != nil {
		r.lokiConfigData.Overrides[tenant.Name] = *tenant.Spec.Limits.Loki
	}
	// update the global loki config
	if r.Config.Spec.Loki.Config != nil {
		r.lokiConfigData.LokiConfigSpec = *r.Config.Spec.Loki.Config
	} else {
		r.lokiConfigData.LokiConfigSpec = observabilityv1alpha1.LokiConfigSpec{}
	}
}

func (r *TenantReconciler) updateTempoConfigmapData(ctx context.Context, tenant *observabilityv1alpha1.Tenant) {
	if tenant.Spec.Limits != nil && tenant.Spec.Limits.Tempo != nil {
		r.tempoConfigData.Overrides[tenant.Name] = *tenant.Spec.Limits.Tempo
	}
}

func (r *TenantReconciler) getMimirConfigMap(ctx context.Context) error {
	existingConfigmap := &corev1.ConfigMap{}

	currentTenantData := mimirConfigData{}

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

func (r *TenantReconciler) getLokiConfigMap(ctx context.Context) error {
	existingConfigmap := &corev1.ConfigMap{}

	currentTenantData := lokiConfigData{}

	err := r.Get(ctx, types.NamespacedName{Name: r.Config.Spec.Loki.ConfigMap.Name, Namespace: r.Config.Spec.Loki.ConfigMap.Namespace}, existingConfigmap)
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
		if tenantData, ok := existingConfigmap.Data[r.Config.Spec.Loki.ConfigMap.Key]; ok {
			yaml.Unmarshal([]byte(tenantData), &currentTenantData)
		} else {
			// TODO: handle error properly
		}
	} else {
		// TODO: handle error properly
	}
	r.lokiConfigData = currentTenantData
	return nil
}

func (r *TenantReconciler) getTempoConfigMap(ctx context.Context) error {
	existingConfigmap := &corev1.ConfigMap{}

	currentTenantData := tempoConfigData{}

	err := r.Get(ctx, types.NamespacedName{Name: r.Config.Spec.Tempo.ConfigMap.Name, Namespace: r.Config.Spec.Tempo.ConfigMap.Namespace}, existingConfigmap)
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
		if tenantData, ok := existingConfigmap.Data[r.Config.Spec.Tempo.ConfigMap.Key]; ok {
			yaml.Unmarshal([]byte(tenantData), &currentTenantData)
		} else {
			// TODO: handle error properly
		}
	} else {
		// TODO: handle error properly
	}
	r.tempoConfigData = currentTenantData
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
			&corev1.ConfigMap{}, // TODO: change this watch to limit to the configmaps we care about
			handler.EnqueueRequestsFromMapFunc(r.findObjectsToReconcile),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		Watches(
			&observabilityv1alpha1.Config{},
			handler.EnqueueRequestsFromMapFunc(r.findObjectsToReconcile),
			builder.WithPredicates(predicate.ResourceVersionChangedPredicate{}),
		).
		Complete(r)
}

func (r *TenantReconciler) findObjectsToReconcile(ctx context.Context, obj client.Object) []reconcile.Request {
	if err := r.Get(ctx, types.NamespacedName{Name: "config"}, r.Config); err != nil {
		if apierrs.IsNotFound(err) {
			// log.Info("Unable to fetch Tenant - skipping", "name", tenantInstance.Name)
			return []reconcile.Request{}
		}
		return []reconcile.Request{}
	}

	if configmap, ok := obj.(*corev1.ConfigMap); ok {
		if (configmap.GetName() == r.Config.Spec.Mimir.ConfigMap.Name && configmap.GetNamespace() == r.Config.Spec.Mimir.ConfigMap.Namespace) ||
			(configmap.GetName() == r.Config.Spec.Loki.ConfigMap.Name && configmap.GetNamespace() == r.Config.Spec.Loki.ConfigMap.Namespace) ||
			(configmap.GetName() == r.Config.Spec.Tempo.ConfigMap.Name && configmap.GetNamespace() == r.Config.Spec.Tempo.ConfigMap.Namespace) {
			tenantList := &observabilityv1alpha1.TenantList{}
			err := r.List(ctx, tenantList, &client.ListOptions{})
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
	} else if _, ok := obj.(*observabilityv1alpha1.Config); ok {
		tenantList := &observabilityv1alpha1.TenantList{}
		err := r.List(ctx, tenantList, &client.ListOptions{})
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
	}
	return []reconcile.Request{}
}

func ignoreNotFound(err error) error {
	if apierrs.IsNotFound(err) {
		return nil
	}
	return err
}
