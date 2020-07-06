/*


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
	"reflect"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	servicemanagerv1 "servicemanager/api/v1"
)

type OwnResource interface {
	// 获取内部资源的实体类
	MakeOwnResource(instance *servicemanagerv1.ServiceManager,logger logr.Logger,scheme *runtime.Scheme)(interface{}, error)

	// 校验资源是否存在
	OwnResourceExist(instance *servicemanagerv1.ServiceManager,client client.Client,logger logr.Logger)(bool, interface{},error)

	// 获取内部资源的状态并修改自定资源的状态
	UpdateOwnerResources(instance *servicemanagerv1.ServiceManager,client client.Client,logger logr.Logger) error

	// 发布内部资源
	ApplyOwnResource(instance *servicemanagerv1.ServiceManager,client client.Client,logger logr.Logger,scheme *runtime.Scheme)error
}



// ServiceManagerReconciler reconciles a ServiceManager object
type ServiceManagerReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=servicemanager.servicemanager.io,resources=servicemanagers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=servicemanager.servicemanager.io,resources=servicemanagers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=statefulSet,verbs=get;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployment,verbs=get;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=service,verbs=get;update;patch;delete
func (r *ServiceManagerReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	logger := r.Log.WithValues("servicemanager", req.NamespacedName)
	serviceManager := &servicemanagerv1.ServiceManager{}

	if err := r.Get(ctx,req.NamespacedName,serviceManager); err != nil{
		logger.Error(err,"获取serviceManager失败！")

		return ctrl.Result{}, err
	}

	// 如果存在，获取own资源
	ownResources,err := r.getOwnResource(serviceManager)
	if err != nil {
		logger.Error(err,"获取ownResource失败！")
	}
	var success = true
	for _,ownResource := range ownResources {
		// 发布或者更新子资源
		if err := ownResource.ApplyOwnResource(serviceManager,r.Client,logger,r.Scheme); err != nil{
			success = false
		}
	}

	// 获取更新内置资源的状态，并且修改自定义资源的crd
	newServiceManager := serviceManager.DeepCopy()
	for _,ownResource := range ownResources {
		// 发布或者更新子资源
		if err := ownResource.UpdateOwnerResources(newServiceManager,r.Client,logger); err != nil{
			success = false
		}
	}

	// 更新newServiceManager
	if newServiceManager != nil && !reflect.DeepEqual(serviceManager.Status,newServiceManager.Status) {
		if err := r.Status().Update(ctx,newServiceManager); err != nil{
			// 这里不处理
			r.Log.Error(err, "unable to update Unit status")
		}

	}

	if !success{
		// 调谐失败
		logger.Info("更新内置资源失败，将监听资源再次放入到workqueue里")
		return ctrl.Result{},err
	}else{
		logger.Info("更新内置资源成功！")
		return ctrl.Result{},nil
	}

	return ctrl.Result{}, nil
}

func (r *ServiceManagerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&servicemanagerv1.ServiceManager{}).
		Complete(r)
}

func (r *ServiceManagerReconciler) getOwnResource(instance *servicemanagerv1.ServiceManager) ([]OwnResource, error) {
	var ownResources []OwnResource
	if instance.Spec.Category == "Deployment" {
		ownDeployment := &servicemanagerv1.OwnDeployment{
			Category:instance.Spec.Category,
		}
		ownResources = append(ownResources, ownDeployment)

	} else {
		// statefulset留着后面写
		/*ownStatefulSet := &servicemanagerv1.OwnStatefulSet{
			Spec: appsv1.StatefulSetSpec{
				Replicas:    instance.Spec.Replicas,
				Selector:    instance.Spec.Selector,
				Template:    instance.Spec.Template,
				ServiceName: instance.Name,
			},
		}

		ownResources = append(ownResources, ownStatefulSet)*/
	}

	if instance.Spec.Port != nil {
		ownService := &servicemanagerv1.OwnService{
			Port:instance.Spec.Port,
		}
		ownResources = append(ownResources, ownService)
	}

	return ownResources,nil

}



