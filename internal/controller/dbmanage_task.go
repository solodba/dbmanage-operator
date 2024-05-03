package controller

import (
	"context"

	operatorcodehorsecomv1beta1 "github.com/solodba/dbmanage-operator/api/v1beta1"
	"k8s.io/apimachinery/pkg/types"
)

// 停止任务循环
func (r *DbManageReconciler) StopLoopTask() {

}

// 开启任务循环
func (r *DbManageReconciler) StartLoopTask() {
	for _, dbManage := range r.DbManageQueue {
		if !dbManage.Spec.Enable {
			dbManage.Status.Active = false
			operatorcodehorsecomv1beta1.L().Info().Msgf("%s任务没有开启", dbManage.Name)
			// 更新状态
			r.UpdateDbManageStatus(dbManage)
			continue
		}

	}
}

// 更新任务的Status信息
func (r *DbManageReconciler) UpdateDbManageStatus(dbManage *operatorcodehorsecomv1beta1.DbManage) {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	namespaceName := types.NamespacedName{
		Namespace: dbManage.Namespace,
		Name:      dbManage.Name,
	}
	ctx := context.TODO()
	dbMangeK8s := operatorcodehorsecomv1beta1.NewDbManage()
	err := r.Client.Get(ctx, namespaceName, dbMangeK8s)
	if err != nil {
		operatorcodehorsecomv1beta1.L().Error().Msgf("获取k8s中DbManage任务%s失败", dbManage.Name)
		return
	}
	dbMangeK8s.Status = dbManage.Status
	err = r.Client.Status().Update(ctx, dbMangeK8s)
	if err != nil {
		operatorcodehorsecomv1beta1.L().Error().Msgf("更新k8s中DbManage任务%s的Status状态失败", dbManage.Name)
		return
	}
}
