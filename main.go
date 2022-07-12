package main

import (
	"flag"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
	"path/filepath"
	"time"
)

type Controller struct {
	indexer cache.Indexer
	queue workqueue.RateLimitingInterface
	informer cache.Controller
}

// Pod控制器
func NewController(queue workqueue.RateLimitingInterface, indexer cache.Indexer, informer cache.Controller)  *Controller{
	return &Controller{
		indexer: indexer,
		queue: queue,
		informer: informer,
	}
}

func (c *Controller) Run(threadiness int, stopCh chan struct{})  {
	defer runtime.HandleCrash()

	// 停止控制器后关掉队列
	defer c.queue.ShuttingDown()

	// 启动控制器da
	klog.Infof("starting pod controller")

	//启动控制器框架
	go c.informer.Run(stopCh)

	// 等待所有相关缓存同步完成，然后再开始处理队列中的数据
	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced){
		runtime.HandleError(fmt.Errorf("time out waiting for caches to sync"))
	}

	// 启动worker处理元素
	for i:= 0; i<threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}
	
	<- stopCh
	klog.Infof("stopped pod controller")
}

func (c *Controller) runWorker() {
	for c.processNextItem(){
	}
}

func (c *Controller) processNextItem() bool {
	// 从workqueue中取出一个元素
	key, quit := c.queue.Get()
	if quit {
		return false
	}

	// 通知队列已经处理了该key
	defer c.queue.Done(key)

	// 根据key去处理我们的业务逻辑
	err := c.syncToStdout(key.(string))
	c.handlerErr(err, key)
	return true
}


func (c *Controller) syncToStdout(key string) error {
	obj, exists,err := c.indexer.GetByKey(key)
	if err != nil {
		klog.Errorf("Forget object with key %s from indexer failed with %v", key, err)
		return nil
	}

	if  !exists {
		fmt.Printf("Pod %s does not exist anymore\n", key)
	}else {
		fmt.Printf("Sync/Add/Update for Pod %s\n", obj.(*v1.Pod).GetName())
	}

	return nil
}

func (c *Controller) handlerErr(err error, key interface{}) {
	if err == nil {
		c.queue.Forget(key)
		return
	}

	if c.queue.NumRequeues(key) < 5 {
		c.queue.AddRateLimited(key)
		return
	}
	c.queue.Forget(key)
	runtime.HandleError(err)

	// 不允许再重试
}


func main() {
	clientset, err := initClient()
	if err != nil {
		klog.Fatal(err)
	}

	podListWatcher := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "pods", v1.NamespaceDefault, fields.Everything())
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter( ))
	indexer, informer := cache.NewIndexerInformer(podListWatcher, &v1.Pod{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			fmt.Println("Add func")
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err != nil {
				queue.Add(key)
			}

		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			fmt.Println("Update func")
			key, err := cache.MetaNamespaceKeyFunc(newObj)
			if err != nil {
				queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Println("Delete func")
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err != nil {
				queue.Add(key)
			}
		},
	}, cache.Indexers{})
	controller := NewController(queue, indexer,informer)

	stopCh := make(chan struct{})
	defer close(stopCh)
	go controller.Run(1, stopCh)

}

func initClient() (*kubernetes.Clientset, error) {
	var err error
	var config *rest.Config
	var kubeconfig *string

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "absolute path to kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "","absolute path to kubeconfig file")
	}

	flag.Parse()

	// 使用sa创建集群配置（InCluster模式）,需要去配置对应的RBAC权限，默认的是default->无权限获取deployments的List权限
	if config, err = rest.InClusterConfig(); err != nil {
		//使用kubeconfig文件创建集群配置
		if config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig); err != nil {
			panic(err.Error())
		}
	}

	return kubernetes.NewForConfig(config)
}
