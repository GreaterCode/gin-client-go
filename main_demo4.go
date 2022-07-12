package main

import (
	"flag"
	"fmt"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"time"
)

// informer示例
func main4() {
	var err error
	var config *rest.Config
	var kubeconfig *string

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "absolute path to kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to kubeconfig file")
	}

	flag.Parse()

	// 使用sa创建集群配置（InCluster模式）,需要去配置对应的RBAC权限，默认的是default->无权限获取deployments的List权限
	if config, err = rest.InClusterConfig(); err != nil {
		//使用kubeconfig文件创建集群配置
		if config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig); err != nil {
			panic(err.Error())
		}
	}

	//创建clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// 初始化informer factory(为了测试方便这里设置每30s重新List一次)
	informerFactory := informers.NewSharedInformerFactory(clientset, time.Second*30)
	// 对Deployment监听
	deployInformer := informerFactory.Apps().V1().Deployments()
	// 创建Informer（相当于注册到工厂中去，这样下面启动的时候就会失去List && Watch对应的资源）
	informer := deployInformer.Informer()
	// 创建Lister
	deployLister := deployInformer.Lister()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    onAdd,
		UpdateFunc: onUpdate,
		DeleteFunc: onDelete,
	})

	stopper := make(chan struct{})
	defer close(stopper)

	//启动informer, Lister && Watch
	informerFactory.Start(stopper)
	// 等待所有启动的Informer的缓存被同步
	informerFactory.WaitForCacheSync(stopper)

	//从本地缓存中获取default中所有的deployment列表
	deployments, err := deployLister.Deployments("default").List(labels.Everything())
	if err != nil {
		panic(err)
	}
	for idx, deployment := range deployments {
		fmt.Printf("%d->%s\n", idx+1, deployment.Name)
	}
	<-stopper
}

func onDelete(obj interface{}) {
	deploy := obj.(*v1.Deployment)
	fmt.Println("delete a deployment:", deploy.Name)
}

func onUpdate(old, new interface{}) {
	oldDeploy := old.(*v1.Deployment)
	newDeploy := new.(*v1.Deployment)
	fmt.Println("update a deployment:", oldDeploy.Name, newDeploy.Name)
}
func onAdd(obj interface{}) {
	deploy := obj.(*v1.Deployment)
	fmt.Println("Add a deployment:", deploy.Name)
}
