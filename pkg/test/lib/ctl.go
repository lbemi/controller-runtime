package lib

import (
	"context"

	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Ctl struct {
}

func (c *Ctl) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	klog.Info(req.NamespacedName)
	return reconcile.Result{}, nil
}
