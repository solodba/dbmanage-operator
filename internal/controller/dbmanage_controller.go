/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	operatorcodehorsecomv1beta1 "github.com/solodba/dbmanage-operator/api/v1beta1"
)

// DbManageReconciler reconciles a DbManage object
type DbManageReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	DbManageQueue map[string]*operatorcodehorsecomv1beta1.DbManage
	Lock          sync.RWMutex
}

//+kubebuilder:rbac:groups=operator.codehorse.com,resources=dbmanages,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=operator.codehorse.com,resources=dbmanages/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=operator.codehorse.com,resources=dbmanages/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DbManage object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.3/pkg/reconcile
func (r *DbManageReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	dbManageK8s := operatorcodehorsecomv1beta1.NewDbManage()
	err := r.Client.Get(ctx, req.NamespacedName, dbManageK8s)
	if err != nil {
		if errors.IsNotFound(err) {
			// 该任务不存在, 则从队列中删除
			r.DeleteQueue(dbManageK8s)
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, err
	}
	// 任务存在, 则对比队列中的任务是否发生改变
	if dbManage, ok := r.DbManageQueue[dbManageK8s.Name]; ok {
		if reflect.DeepEqual(dbManageK8s, dbManage) {
			return ctrl.Result{}, fmt.Errorf("%s任务没有发生任务变化", dbManageK8s.Name)
		}
	}
	// 任务发生变化或者没有该任务则添加进队列
	r.AddQueue(dbManageK8s)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DbManageReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&operatorcodehorsecomv1beta1.DbManage{}).
		Complete(r)
}
