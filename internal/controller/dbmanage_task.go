package controller

import (
	"context"
	"fmt"
	"time"

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
		// 获取任务开始时间
		taskDelay := r.GetTaskDelaySeconds(dbManage.Spec.StartTime)
		if taskDelay.Hours() < 1 {
			operatorcodehorsecomv1beta1.L().Info().Msgf("%s任务还有%.1f分钟后开始执行", dbManage.Name, taskDelay.Minutes())
		} else {
			operatorcodehorsecomv1beta1.L().Info().Msgf("%s任务还有%.1f小时后开始执行", dbManage.Name, taskDelay.Hours())
		}
		// 更新任务状态
		dbManage.Status.Active = true
		dbManage.Status.NextTime = int(r.GetTaskNextTime(taskDelay.Seconds()).Unix())
		r.UpdateDbManageStatus(dbManage)
		// 执行指定任务
		ticker := time.NewTicker(taskDelay)
		r.Tickers = append(r.Tickers, ticker)
		r.Wg.Add(1)
		go func(dbManage *operatorcodehorsecomv1beta1.DbManage) {
			defer r.Wg.Done()
			for {
				<-ticker.C
				// 重置ticker
				ticker.Reset(time.Minute * time.Duration(dbManage.Spec.Period))
				operatorcodehorsecomv1beta1.L().Info().Msgf("%s任务将在%d分钟后循环执行", dbManage.Name, dbManage.Spec.Period)
				// 更新任务状态
				dbManage.Status.Active = true
				dbManage.Status.NextTime = int(r.GetTaskNextTime(float64(dbManage.Spec.Period * 60)).Unix())
				switch dbManage.Spec.Flag {
				case 0:
					// 备份任务
					err := r.DbBackupTask(dbManage)
					if err != nil {
						operatorcodehorsecomv1beta1.L().Error().Msgf("%s备份任务失败, 原因: %s", dbManage.Name, err.Error())
						dbManage.Status.LastTaskResult = fmt.Sprintf("%s备份任务失败, 原因: %s", dbManage.Name, err.Error())
					}
					operatorcodehorsecomv1beta1.L().Error().Msgf("%s备份任务成功", dbManage.Name)
					dbManage.Status.LastTaskResult = fmt.Sprintf("%s备份任务成功", dbManage.Name)
				case 1:
					// 巡检任务
					// TODO

				default:
				}
				// 更新任务状态
				r.UpdateDbManageStatus(dbManage)
			}
		}(dbManage)
	}
	r.Wg.Wait()
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
		operatorcodehorsecomv1beta1.L().Error().Msgf("获取k8s中DbManage任务%s失败, 原因:%s", dbManage.Name, err.Error())
		return
	}
	dbMangeK8s.Status = dbManage.Status
	err = r.Client.Status().Update(ctx, dbMangeK8s)
	if err != nil {
		operatorcodehorsecomv1beta1.L().Error().Msgf("更新k8s中DbManage任务%s的Status状态失败, 原因:%s", dbManage.Name, err.Error())
		return
	}
}
