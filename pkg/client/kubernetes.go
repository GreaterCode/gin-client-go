package client

import (
	"errors"
	"flag"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog"
	"path/filepath"
	"sync"
)

var onceClient sync.Once
var onceConfig sync.Once
var Kubeconfig *rest.Config
var KubeClientSet *kubernetes.Clientset

func GetK8sClientSet() (*kubernetes.Clientset, error) {
	onceClient.Do(func() {
		config, err := GetRestConfig()
		if err != nil {
			return
		}
		KubeClientSet, err = kubernetes.NewForConfig(config)
		if err != nil {
			klog.Fatal(err)
			return
		}
	})
	return KubeClientSet, nil
}

func GetRestConfig() (config *rest.Config, err error) {
	onceConfig.Do(func() {
		var kubeconfig *string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "")
		} else {
			klog.Fatal("read config error, config is empty")
			err = errors.New("read config error, config is empty")
			return
		}
		flag.Parse()
		Kubeconfig, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			klog.Fatal(err)
			return
		}
		return
	})
	return Kubeconfig, nil
}
