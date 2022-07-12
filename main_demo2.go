package main

import (
	"context"
	"flag"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

// 获取defaults命名空间下所有deployment列表
func main1() {
	var err error
	var config *rest.Config
	var kubeconfig *string
	ctx := context.Background()

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

	//创建clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	//使用clientset获取deployments
	deployments, err := clientset.AppsV1().Deployments("kube-system").List(ctx, metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	for idx, item := range deployments.Items {
		fmt.Printf("%d->%s\n", idx+1, item.Name)
	}

	pods, err:= clientset.CoreV1().Pods("kube-system").List(ctx, metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for idx, item := range pods.Items {
		fmt.Printf("%d->%s\n", idx+1, item.Name)
	}

}
