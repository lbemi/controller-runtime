package main

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/klog/v2"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func main() {

	s := runtime.NewScheme()
	//AddToScheme将此clientset的所有类型添加到新的scheme中
	scheme.AddToScheme(s)
	//corev1.AddToScheme(s) // 只将corev1中的gvk添加到新的scheme中
	//p := &corev1.Pod{}
	//fmt.Println("scheme中的kind是通过reflect.TypeOf(p).Elem().Name获取结构体的名字", reflect.TypeOf(p).Elem().Name())
	//fmt.Println(s)
	objectKinds, _, _ := s.ObjectKinds(&corev1.Pod{})
	fmt.Println("获取gvk", objectKinds)
	//return
	//创建manager
	mgr, err := manager.New(kubernetesConfig(), manager.Options{
		Logger: logf.Log.WithName("test"),
	})
	if err != nil {
		klog.Fatalf(" make manager failed : %s", err.Error())
	}
	ctx := context.Background()
	kinds, _, err := mgr.GetScheme().ObjectKinds(&corev1.Pod{})
	if err != nil {
		klog.Fatalf("get kinds error : %s", err.Error())
	}
	klog.Info("kinds", kinds)
	go func() {
		time.Sleep(3 * time.Second)
		pod := &corev1.Pod{}
		err = mgr.GetClient().Get(ctx, types.NamespacedName{
			Namespace: "default",
			Name:      "nginx-7f65fcf556-w8q9s",
		}, pod)
		if err != nil {
			klog.Fatalf("get pod error : %s", err.Error())
		}
		klog.Info("pod: ", pod.Name)
		//cache使用的informer,最终在mgr.Start()函数中启动
		//podInformer, _ := mgr.GetCache().GetInformer(ctx, &corev1.Pod{})
		//podInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		//	AddFunc: func(obj interface{}) {
		//		fmt.Println("add")
		//	},
		//	UpdateFunc: func(oldObj, newObj interface{}) {
		//		fmt.Println("update")
		//	},
		//	DeleteFunc: func(obj interface{}) {
		//		fmt.Println("delete")
		//	},
		//})
	}()

	err = mgr.Start(ctx)
	if err != nil {
		klog.Fatalf("start manager error : %s", err.Error())
	}
}

// pkg/test/test.go
