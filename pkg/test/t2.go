package main

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/source"
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
		Reconciler: &lib.Ctl{},
	})
	if err != nil {
		klog.Fatalf(" make controller failed : %s", err.Error())
	}

	//手动添加添加watch pod资源
	h := &handler.TypedEnqueueRequestForObject[*corev1.Pod]{}
	err = ctl.Watch(
		source.Kind(mgr.GetCache(), &corev1.Pod{}, h),
	)
	if err != nil {
		klog.Fatalf(" add watch failed : %s", err.Error())
	}

	//交给manager一起启动
	err = mgr.Add(ctl)
	if err != nil {
		klog.Fatalf(" add controller failed : %s", err.Error())
	}

	ctx := context.Background()
	mgr.Add(lib.NewWeb(h, ctl))

	err = mgr.Start(ctx)
	if err != nil {
		klog.Fatalf("start manager error : %s", err.Error())
	}
}
