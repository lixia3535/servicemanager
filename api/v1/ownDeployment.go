package v1

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type OwnDeployment struct {
	Category string
}


func (ownDeployment *OwnDeployment) MakeOwnResource(instance *ServiceManager,logger logr.Logger,scheme *runtime.Scheme) (interface{}, error){
	var label  = map[string]string{
		"app": instance.Name,
	}
	var selector = &metav1.LabelSelector{
		MatchLabels: label,
	}
	deployment := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:                       instance.Name,
			Namespace:                  instance.Namespace,
		},
		Spec:v1.DeploymentSpec{
			Replicas:                instance.Spec.Replicas,
			Template:                instance.Spec.Template,
			Selector:				 selector,
		},
	}

	deployment.Spec.Template.Labels = label
	if err :=controllerutil.SetControllerReference(instance,deployment,scheme); err != nil{
		msg := fmt.Sprintf("set controllerReference for service %s/%s failed", instance.Namespace, instance.Name)
		logger.Error(err, msg)
		return nil, err
	}
	return deployment,nil
}

// 校验资源是否存在
func (ownDeployment *OwnDeployment) OwnResourceExist(instance *ServiceManager,client client.Client,logger logr.Logger) (bool, interface{},error){

	deployment := &v1.Deployment{}
	// 查看k8s集群中是否存在service资源
	if err := client.Get(context.Background(),types.NamespacedName{Name:instance.Name,Namespace:instance.Namespace},deployment); err != nil{
		return false,nil,err
	}

	return true,deployment,nil
}

// 获取内部资源的状态并修改自定资源的状态
func (ownDeployment *OwnDeployment) UpdateOwnerResources(instance *ServiceManager,client client.Client,logger logr.Logger) error{

	deployment := &v1.Deployment{}
	if err := client.Get(context.Background(),types.NamespacedName{Name:instance.Name,Namespace:instance.Namespace},deployment); err != nil{
		logger.Error(err,"service 资源不存在！")
		return err
	}

	instance.Status.LastUpdateTime = metav1.Now()
	instance.Status.DeploymentStatus = deployment.Status

	return nil
}

// 发布内部资源
func (ownDeployment *OwnDeployment) ApplyOwnResource(instance *ServiceManager,client client.Client,logger logr.Logger,scheme *runtime.Scheme) error{

	// 首先查看资源是否存在
	exsit,found,err := ownDeployment.OwnResourceExist(instance,client,logger)
	if err != nil {
		logger.Error(err,"deployment 资源不存在！")
		// return err
	}

	deployment,err := ownDeployment.MakeOwnResource(instance,logger,scheme)
	newDeployment,ok := deployment.(*v1.Deployment)
	if !ok {
		logger.Error(err,"deployment 结构体转化失败！")
		return err
	}
	if err != nil {
		logger.Error(err,"获取deployment资源失败！")
		return err
	}
	if exsit {
		// 更新
		founDeployment,ok := found.(*v1.Deployment)
		if  ! ok{
			logger.Error(err,"deployment 结构体转化失败！")
			return err
		}

		// 如果不等，则更新Service资源
		if founDeployment != nil && !reflect.DeepEqual(founDeployment.Spec,newDeployment.Spec) {
			err := client.Update(context.Background(),newDeployment)
			if err != nil {
				logger.Error(err,"deployment更新失败！")
				return err
			}
		}
	}else{
		// 创建
		err := client.Create(context.Background(),newDeployment)
		if err != nil {
			logger.Error(err,"deployment 创建失败！")
			return err
		}
	}
	return nil

}
