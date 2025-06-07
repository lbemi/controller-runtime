package lib

import (
	"context"

	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	cc "sigs.k8s.io/controller-runtime/pkg/internal/controller"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Web struct {
	E   handler.TypedEventHandler[*v1.Pod, reconcile.Request]
	Ctl controller.Controller
}

func NewWeb(e handler.TypedEventHandler[*v1.Pod, reconcile.Request], ctl controller.Controller) *Web {
	return &Web{E: e, Ctl: ctl}
}

func (w *Web) Start(context.Context) error {
	r := gin.New()
	r.GET("/add", func(c *gin.Context) {
		pod := &v1.Pod{}
		pod.Name = "test-pod"
		pod.Namespace = "default"
		w.E.Create(
			context.Background(), event.TypedCreateEvent[*v1.Pod]{Object: pod},
			w.Ctl.(*cc.Controller[reconcile.Request]).Queue,
		)
		c.JSON(200, gin.H{"message": "ok"})
	})

	return r.Run(":8081")
}
