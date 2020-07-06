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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ServiceManagerSpec defines the desired state of ServiceManager


type ServiceManagerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of ServiceManager. Edit ServiceManager_types.go to remove/update
	// Category 只有两种可能 deployment statefulset
	// 这个注释表示该字段的值只能是Deployment 或者 Statefulset
	// +kubebuilder:validation:Enum=Deployment;Statefulset
	Category string `json:"category,omitempty"`
	
	// 标签选择器
	Selector map[string]string `json:"selector,omitempty"`
	
	// 引用的statefulset deployment的template
	Template corev1.PodTemplateSpec `json:"template,omitempty"`

	// 副本数 最大不超过10
	// +kubebuilder:validation:Maximum=10
	Replicas *int32 `json:"replicas,omitempty"`
	
	//端口号 端口号做大超过65535 服务端口号
	// +kubebuilder:validation:Maximum=65535
	Port *int32 `json:"port,omitempty"`

	// +kubebuilder:validation:Maximum=65535
	Targetport *int32 `json:"targetport,omitempty"`
	
}

// ServiceManagerStatus defines the observed state of ServiceManager
type ServiceManagerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Replicas int32 `json:"replicas,omitempty"`
	LastUpdateTime metav1.Time `json:"last_update_time,omitempty"`
	DeploymentStatus appsv1.DeploymentStatus `json:"deployment_status,omitempty"`
	ServiceStatus corev1.ServiceStatus `json:"service_status,omitempty"`
}
// 这里，Spec和Status均是ServiceManager的成员变量，Status并不像Pod.Status一样，是Pod的subResource.因此，
// 如果我们在controller的代码中调用到Status().Update(),会触发panic，
// 并报错：the server could not find the requested resource
// 如果我们想像k8s中的设计那样，那么就要遵循k8s中status subresource的使用规范：
// kubebuilder:subresource:status
// 用户只能指定一个CRD实例的spec部分；
// CRD实例的status部分由控制器进行变更。

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:subresource:scale:selectorpath=.spec.selector,specpath=.spec.replicas,statuspath=.status.replicas
// ServiceManager is the Schema for the servicemanagers API
type ServiceManager struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServiceManagerSpec   `json:"spec,omitempty"`
	Status ServiceManagerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ServiceManagerList contains a list of ServiceManager
type ServiceManagerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ServiceManager `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ServiceManager{}, &ServiceManagerList{})
}
