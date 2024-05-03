package controller

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	operatorcodehorsecomv1beta1 "github.com/solodba/dbmanage-operator/api/v1beta1"
)

// 创建目录
func (r *DbManageReconciler) CreateDir(dbMange *operatorcodehorsecomv1beta1.DbManage) error {
	switch dbMange.Spec.Flag {
	case 0:
		return os.MkdirAll("/tmp/dbbackup", 0755)
	case 1:
		return os.MkdirAll("/tmp/dbcheck", 0755)
	default:
		return fmt.Errorf("不支持该类型的任务")
	}
}

// 数据库备份任务
func (r *DbManageReconciler) DbBackupTask(dbMange *operatorcodehorsecomv1beta1.DbManage) error {
	err := r.CreateDir(dbMange)
	if err != nil {
		return err
	}
	backupCmd := fmt.Sprintf("mysqldump -u%s -p%s -h%s -P%d --all-databases > /tmp/dbbackup/%s.sql",
		dbMange.Spec.Origin.Username,
		dbMange.Spec.Origin.Password,
		dbMange.Spec.Origin.Host,
		dbMange.Spec.Origin.Port,
		time.Now().Format("0102150405"))
	_, err = exec.Command(backupCmd).Output()
	if err != nil {
		return err
	}
	return nil
}
