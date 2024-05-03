package controller

import (
	operatorcodehorsecomv1beta1 "github.com/solodba/dbmanage-operator/api/v1beta1"
)

// 任务添加到队列
func (r *DbManageReconciler) AddQueue(dbManage *operatorcodehorsecomv1beta1.DbManage) {
	if r.DbManageQueue == nil {
		r.DbManageQueue = make(map[string]*operatorcodehorsecomv1beta1.DbManage)
	}
	r.DbManageQueue[dbManage.Name] = dbManage
	// 停止和开启任务循环
	r.StopLoopTask()
	go r.StartLoopTask()
}

// 任务从队列中删除
func (r *DbManageReconciler) DeleteQueue(dbManage *operatorcodehorsecomv1beta1.DbManage) {
	delete(r.DbManageQueue, dbManage.Name)
	// 停止和开启任务循环
	r.StopLoopTask()
	go r.StartLoopTask()
}
