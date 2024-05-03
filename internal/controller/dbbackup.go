package controller

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/minio/minio-go/v7"
	operatorcodehorsecomv1beta1 "github.com/solodba/dbmanage-operator/api/v1beta1"
)

// 数据库备份任务
func (r *DbManageReconciler) DbBackupTask(dbManage *operatorcodehorsecomv1beta1.DbManage) error {
	// 备份任务
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
	fileName := time.Now().Format("0102150405") + ".sql"
	backupCmd := fmt.Sprintf("mysqldump -u%s -p%s -h%s -P%d --all-databases > %s/%s",
		dbManage.Spec.Origin.Username,
		dbManage.Spec.Origin.Password,
		dbManage.Spec.Origin.Host,
		dbManage.Spec.Origin.Port,
		dbBackupDir,
		fileName)
	_, err = exec.Command("bash", "-c", backupCmd).Output()
	if err != nil {
		return err
	}
	// 同步备份文件到MinIO
	minioClient, err := r.InitialMinioClient(dbManage)
	if err != nil {
		return err
	}
	f, err := os.Open(fmt.Sprintf("%s/%s", dbBackupDir, fileName))
	if err != nil {
		operatorcodehorsecomv1beta1.L().Error().Msgf("打开备份文件%s失败", fileName)
		return err
	}
	_, err = minioClient.PutObject(
		context.TODO(),
		dbManage.Spec.Destination.BucketName,
		fmt.Sprintf("%s/%s", dbBackupDir, fileName),
		f,
		-1,
		minio.PutObjectOptions{})
	if err != nil {
		operatorcodehorsecomv1beta1.L().Error().Msgf("上传备份文件%s到minio失败", fileName)
		return err
	}
	return nil
}
