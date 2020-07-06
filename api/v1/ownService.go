package v1

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type OwnService struct {
	Port *int32
}
// 获取内部资源的实体类
func (ownService *OwnService) MakeOwnResource(instance *ServiceManager,logger logr.Logger,scheme *runtime.Scheme) (interface{}, error){
	var label  = map[string]string{
		"app": instance.Name,
	}
	objectMeta := metav1.ObjectMeta{
		Name:instance.Name,
		Namespace:instance.Namespace,

	}
	servicePort := []corev1.ServicePort{
		corev1.ServicePort{
			TargetPort: intstr.IntOrString{intstr.Int,*instance.Spec.Targetport,""},
			NodePort:   *instance.Spec.Port,
			Port:*instance.Spec.Port,
		},
	}
	serviceSpec := corev1.ServiceSpec{
		Selector:label,
		Type:corev1.ServiceTypeNodePort,
		Ports:servicePort,

	}
	service := &corev1.Service{
		ObjectMeta: objectMeta,
		Spec:serviceSpec,
	}
	if err :=controllerutil.SetControllerReference(instance,service,scheme); err != nil{
		msg := fmt.Sprintf("set controllerReference for service %s/%s failed", instance.Namespace, instance.Name)
		logger.Error(err, msg)
		return nil, err
	}
	return service,nil
}

// 校验资源是否存在
func (ownService *OwnService) OwnResourceExist(instance *ServiceManager,client client.Client,logger logr.Logger) (bool, interface{},error){

	service := &corev1.Service{}
	// 查看k8s集群中是否存在service资源
	if err := client.Get(context.Background(),types.NamespacedName{Name:instance.Name,Namespace:instance.Namespace},service); err != nil{
		return false,nil,err
	}

	return true,service,nil
}

// 获取内部资源的状态并修改自定资源的状态
func (ownService *OwnService) UpdateOwnerResources(instance *ServiceManager,client client.Client,logger logr.Logger) error{

	service := &corev1.Service{}
	if err := client.Get(context.Background(),types.NamespacedName{Name:instance.Name,Namespace:instance.Namespace},service); err != nil{
		logger.Error(err,"service 资源不存在！")
		return err
	}

	instance.Status.LastUpdateTime = metav1.Now()
	instance.Status.ServiceStatus = service.Status

	return nil
}

// 发布内部资源
func (ownService *OwnService) ApplyOwnResource(instance *ServiceManager,client client.Client,logger logr.Logger,scheme *runtime.Scheme) error{

	// 首先查看资源是否存在
	exsit,found,err := ownService.OwnResourceExist(instance,client,logger)
	if err != nil {
		logger.Error(err,"service 资源不存在！")
		// return err
	}

	service,err := ownService.MakeOwnResource(instance,logger,scheme)
	newService,ok := service.(*corev1.Service)
	if !ok {
		logger.Error(err,"service 结构体转化失败！")
		return err
	}
	if err != nil {
		logger.Error(err,"获取service资源失败！")
		return err
	}
	if exsit {
		// 更新
		founService,ok := found.(*corev1.Service)
		if  ! ok{
			logger.Error(err,"service 结构体转化失败！")
			return err
		}
		// 这里有个坑，svc在创建前可能未指定clusterIP，那么svc创建后，
		// 会自动指定clusterIP并修改spec.clusterIP字段，
		// 因此这里要补上。SessionAffinity同理
		newService.Spec.ClusterIP = founService.Spec.ClusterIP
		newService.Spec.SessionAffinity = founService.Spec.SessionAffinity
		newService.ObjectMeta.ResourceVersion = founService.ObjectMeta.ResourceVersion
		// 如果不等，则更新Service资源
		if founService != nil && !reflect.DeepEqual(founService.Spec,newService.Spec) {
			err := client.Update(context.Background(),newService)
			if err != nil {
				logger.Error(err,"service更新失败！")
				return err
			}
		}
	}else{
		// 创建
		err := client.Create(context.Background(),newService)
		if err != nil {
			logger.Error(err,"service 创建失败！")
			return err
		}
	}
	return nil

}
