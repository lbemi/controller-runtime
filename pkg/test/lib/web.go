package lib

import (
	"context"

	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	cc "sigs.k8s.io/controller-runtime/pkg/internal/controller"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Web struct {
	E      handler.TypedEventHandler[client.Object, reconcile.Request]
	Ctl    controller.Controller
	newObj metav1.Object
	scheme *runtime.Scheme
}

func NewWeb(e handler.TypedEventHandler[client.Object, reconcile.Request], ctl controller.Controller, newObj metav1.Object, scheme *runtime.Scheme) *Web {
	return &Web{E: e, Ctl: ctl, newObj: newObj, scheme: scheme}
}

func (w *Web) Start(context.Context) error {
	r := gin.New()

	r.GET("/add", func(c *gin.Context) {
		cm := &v1.ConfigMap{}
		cm.Name = "test-cm"
		cm.Namespace = "default"

		controllerutil.SetOwnerReference(w.newObj, cm, w.scheme) // 设置owner
		if _, ok := w.Ctl.(*cc.Controller[reconcile.Request]); !ok {
			c.JSON(200, gin.H{"message": "error"})
		}
		// 手动出发  reconcile,添加事件
		w.E.Create(
			context.Background(),
			event.TypedCreateEvent[client.Object]{Object: cm},
			w.Ctl.(*cc.Controller[reconcile.Request]).Queue,
		)
		c.JSON(200, gin.H{"message": "ok"})
	})

	return r.Run(":8081")
}
