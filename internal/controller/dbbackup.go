package controller

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	operatorcodehorsecomv1beta1 "github.com/solodba/dbmanage-operator/api/v1beta1"
)

// 数据库备份任务
func (r *DbManageReconciler) DbBackupTask(dbManage *operatorcodehorsecomv1beta1.DbManage) error {
	dbBackupDir := "/tmp/dbbackup"
	_, err := os.Stat(dbBackupDir)
	if err != nil {
		if errx := os.MkdirAll(dbBackupDir, 0700); errx == nil {
			operatorcodehorsecomv1beta1.L().Info().Msgf("%s文件夹创建成功", dbBackupDir)
		} else {
			operatorcodehorsecomv1beta1.L().Info().Msgf("%s文件夹创建失败, 原因: %s", dbBackupDir, errx.Error())
			return errx
		}
	}
	backupCmd := fmt.Sprintf("mysqldump -u%s -p%s -h%s -P%d --all-databases > /tmp/dbbackup/%s.sql",
		dbManage.Spec.Origin.Username,
		dbManage.Spec.Origin.Password,
		dbManage.Spec.Origin.Host,
		dbManage.Spec.Origin.Port,
		time.Now().Format("0102150405"))
	_, err = exec.Command("bash", "-c", backupCmd).Output()
	if err != nil {
		return err
	}
	return nil
}
