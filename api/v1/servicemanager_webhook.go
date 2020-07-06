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

package v1

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var servicemanagerlog = logf.Log.WithName("servicemanager-resource")

func (r *ServiceManager) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-servicemanager-servicemanager-io-v1-servicemanager,mutating=true,failurePolicy=fail,groups=servicemanager.servicemanager.io,resources=servicemanagers,verbs=create;update,versions=v1,name=mservicemanager.kb.io

var _ webhook.Defaulter = &ServiceManager{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
// 这个方法时可以修改结构体 比如添加一些默认参数什么的
func (r *ServiceManager) Default() {
	servicemanagerlog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// +kubebuilder:webhook:verbs=create;update,path=/validate-servicemanager-servicemanager-io-v1-servicemanager,mutating=false,failurePolicy=fail,groups=servicemanager.servicemanager.io,resources=servicemanagers,versions=v1,name=vservicemanager.kb.io

var _ webhook.Validator = &ServiceManager{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
// 下面的方法主要时做校验使用
func (r *ServiceManager) ValidateCreate() error {
	servicemanagerlog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *ServiceManager) ValidateUpdate(old runtime.Object) error {
	servicemanagerlog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *ServiceManager) ValidateDelete() error {
	servicemanagerlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
