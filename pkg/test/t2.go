package main

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/test/lib"
)

func config() *rest.Config {
	homeDir := homedir.HomeDir()
	config, err := clientcmd.BuildConfigFromFlags("", homeDir+"/.kube/config")
	if err != nil {
		klog.Fatalf("Get kubeconfig error : %s", err.Error())
	}
	return config
}

func main() {

	//创建manager
	mgr, err := manager.New(config(), manager.Options{
		Logger: logf.Log.WithName("test"),
	})
	if err != nil {
		klog.Fatalf(" make manager failed : %s", err.Error())
	}
	//创建控制器
	ctl, err := controller.New("test", mgr, controller.Options{
		Reconciler:              &lib.Ctl{},
		MaxConcurrentReconciles: 1, //并发数
	})
	if err != nil {
		klog.Fatalf(" make controller failed : %s", err.Error())
		return
	}

	//手动添加添加watch pod资源
	// h := &handler.TypedEnqueueRequestForObject[*corev1.Pod]{}
	// err = ctl.Watch(
	// 	source.Kind(mgr.GetCache(), &corev1.Pod{}, h2),
	// )

	// if err != nil {
	// 	klog.Fatalf(" add watch failed : %s", err.Error())
	// }

	h := handler.TypedEnqueueRequestForOwner[client.Object]( //设置owner资源监听
		mgr.GetScheme(),
		mgr.GetRESTMapper(),
		&corev1.Pod{},
	)
	// 手动出发 reconcile
	err = mgr.Add(lib.NewWeb(h, ctl, &corev1.Pod{}, mgr.GetScheme()))
	if err != nil {
		return
	}

	ctx := context.Background()
	err = mgr.Start(ctx)
	if err != nil {
		klog.Fatalf("start manager error : %s", err.Error())
	}
}
